package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileMetadata struct {
	OriginalName   string    `json:"original_name"`
	NewName        string    `json:"new_name"`
	MovedAt        time.Time `json:"moved_at"`
	OriginalSize   int64     `json:"original_size"`
	CompressedSize int64     `json:"compressed_size"`
	ItemType       string    `json:"item_type"`
}

func SaveMetadata(targetDir string, newEntries []FileMetadata) error {
	if len(newEntries) == 0 {
		return nil
	}
	path := filepath.Join(targetDir, "metadata.json")
	var existingEntries []FileMetadata

	// Read existing if exists
	if data, err := os.ReadFile(path); err != nil && len(data) > 0 {
		json.Unmarshal(data, &existingEntries)
	}

	// Merge old and new
	allEntries := append(existingEntries, newEntries...)

	// Write
	newData, err := json.MarshalIndent(allEntries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, newData, 0644)
}

func LoadAllMetadata(targetDir string) ([]FileMetadata, error) {
	path := filepath.Join(targetDir, "metadata.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []FileMetadata
	err = json.Unmarshal(data, &entries)
	return entries, err
}

func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
