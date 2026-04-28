package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Tidy(conf *Config) error {
	desktop, _ := os.UserHomeDir()
	desktop = filepath.Join(desktop, "Desktop")

	// Ensure target dir exists
	os.MkdirAll(conf.TargetDir, 0755)

	items, err := os.ReadDir(desktop)
	if err != nil {
		return err
	}

	for _, item := range items {
		// Skip excluded files
		if shouldExclude(item.Name(), conf.ExcludeFiles) {
			continue
		}

		src := filepath.Join(desktop, item.Name())
		info, err := item.Info()
		if err != nil {
			continue
		}

		// File type
		itemType := "file"
		if item.IsDir() {
			itemType = "folder"
		} else if isShortcut(src, item.Name()) {
			itemType = "shortcut"
		}

		// Destination
		timestamp := time.Now().Format("20060102_150405")
		newName := fmt.Sprintf("%s_%s", timestamp, item.Name())
		dest := filepath.Join(conf.TargetDir, newName)

		// Move
		fmt.Printf("Moving %s: %s -> %s\n", itemType, item.Name(), newName)
		if err := os.Rename(src, dest); err != nil {
			fmt.Printf("Error moving %s: %v\n", item.Name(), err)
			continue
		}

		// Metadata
		UpdateMetadata(conf.TargetDir, FileMetadata{
			OriginalName: item.Name(),
			NewName:      newName,
			MovedAt:      time.Now(),
			Size:         info.Size(),
			ItemType:     itemType,
		})
	}
	return nil
}

func shouldExclude(name string, excluded []string) bool {
	for _, ex := range excluded {
		if name == ex {
			return true
		}
	}
	return false
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
