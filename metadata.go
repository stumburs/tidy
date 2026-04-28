package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type FileMetadata struct {
	OriginalName string    `json:"original_name"`
	NewName      string    `json:"new_name"`
	MovedAt      time.Time `json:"moved_at"`
	Size         int64     `json:"size"`
	ItemType     string    `json:"item_type"`
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
