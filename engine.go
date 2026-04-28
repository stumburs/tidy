package main

import (
	"archive/zip"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func Tidy(conf *Config) error {
	home, _ := os.UserHomeDir()
	desktop := filepath.Join(home, "Desktop")

	// Ensure target dir exists
	os.MkdirAll(conf.TargetDir, 0755)

	items, err := os.ReadDir(desktop)
	if err != nil {
		return err
	}

	var batchMetadata []FileMetadata

	for _, item := range items {
		// Skip excluded files
		if shouldExclude(item.Name(), conf.ExcludeFiles) {
			continue
		}

		src := filepath.Join(desktop, item.Name())
		info, _ := item.Info()

		// Destination
		timestamp := time.Now().Format("20060102_150405")
		newName := fmt.Sprintf("%s_%s.zip", timestamp, item.Name())
		dest := filepath.Join(conf.TargetDir, newName)

		bar := progressbar.DefaultBytes(
			info.Size(),
			"Compressing: "+item.Name(),
		)

		err := zipSourceWithProgress(src, dest, bar, conf.CompressionLevel)
		if err != nil {
			fmt.Printf("\n  Failed to compress %s: %v\n", item.Name(), err)
			continue
		}

		zipInfo, _ := os.Stat(dest)

		// Remove original file
		os.RemoveAll(src)

		// Metadata
		batchMetadata = append(batchMetadata, FileMetadata{
			OriginalName:   item.Name(),
			NewName:        newName,
			MovedAt:        time.Now(),
			OriginalSize:   info.Size(),
			CompressedSize: zipInfo.Size(),
			ItemType:       getItemType(item),
		})
		fmt.Println()
	}
	return SaveMetadata(conf.TargetDir, batchMetadata)
}

func shouldExclude(name string, excluded []string) bool {
	return slices.Contains(excluded, name)
}

func isShortcut(path string, name string) bool {
	// Windows .lnk files
	if strings.HasSuffix(strings.ToLower(name), ".lnk") {
		return true
	}

	// Unix Symlinks
	info, err := os.Lstat(path)
	if err == nil && info.Mode()&os.ModeSymlink != 0 {
		return true
	}

	return false
}

func ShowStats(conf *Config) {
	entries, err := LoadAllMetadata(conf.TargetDir)
	if err != nil {
		fmt.Println("Could not read metadata history.")
		return
	}

	var totalOrig, totalComp int64

	for _, e := range entries {
		totalOrig += e.OriginalSize
		totalComp += e.CompressedSize
	}

	saved := totalOrig - totalComp
	// Prevent division by zero
	ratio := 0.0
	if totalOrig > 0 {
		ratio = (float64(saved) / float64(totalOrig)) * 100
	}

	fmt.Println("\n--- tidy stats ---")
	fmt.Printf("Total Items Archived:    %d\n", len(entries))
	fmt.Printf("Original Size Archived:  %s\n", FormatSize(totalOrig))
	fmt.Printf("Compressed Size:         %s\n", FormatSize(totalComp))
	fmt.Printf("Space reclaimed:         %s (%.1f%% reduction)\n", FormatSize(saved), ratio)
	fmt.Println("--------------------")
}

func SearchFiles(conf *Config, query string) {
	entries, err := LoadAllMetadata(conf.TargetDir)
	if err != nil {
		fmt.Println("No archive history found.")
		return
	}
	fmt.Printf("Searching for: '%s'\n", query)
	found := false
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.OriginalName), strings.ToLower(query)) {
			found = true
			fmt.Printf("\nFound: %s\n", e.OriginalName)
			fmt.Printf("  Current Name: %s\n", e.NewName)
			fmt.Printf("  Location:     %s\n", filepath.Join(conf.TargetDir, e.NewName))
			fmt.Printf("  Archived On:  %s\n", e.MovedAt.Format("Jan 02, 2006"))
		}
	}

	if !found {
		fmt.Println("No matches found in your archive.")
	}
}

func zipSourceWithProgress(source, target string, bar *progressbar.ProgressBar, level int) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	writer.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, level)
	})

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, _ = filepath.Rel(filepath.Dir(source), path)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		w, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(io.MultiWriter(w, bar), file)
		return err
	})
}

func getItemType(item os.DirEntry) string {
	if item.IsDir() {
		return "folder"
	}
	return "file"
}
