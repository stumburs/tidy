package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	TargetDir        string   `toml:"target_dir"`
	ExcludeFiles     []string `toml:"exclude_files"`
	CompressionLevel int      `toml:"compression_level"` // 1 (fast) - 9 (slow)
}

func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config
		home, _ := os.UserHomeDir()
		defaultConfig := Config{
			TargetDir:        filepath.Join(home, "Documents", "TidyArchive"),
			ExcludeFiles:     []string{".DS_Store", "desktop.ini"},
			CompressionLevel: 9, // max compression
		}

		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		toml.NewEncoder(f).Encode(defaultConfig)
		return &defaultConfig, nil
	}

	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
