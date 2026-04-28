package main

import (
	"encoding/json"
	"os"
	"time"
)

type FileMetadata struct {
	OriginalName string    `json:"original_name"`
	NewName      string    `json:"new_name"`
	MovedAt      time.Time `json:"moved_at"`
	Size         int64     `json:"size"`
	ItemType     string    `json:"item_type"`
}

func UpdateMetadata(targetDir string, entry FileMetadata) error {
	path := targetDir + "/metadata.json"
	var entries []FileMetadata

	// Read existing if exists
	if data, err := os.ReadFile(path); err != nil {
		json.Unmarshal(data, &entries)
	}

	entries = append(entries, entry)

	newData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, newData, 0644)
}
