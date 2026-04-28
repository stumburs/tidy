package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	configPath := "tidy.config.toml"

	// Check if config exists
	_, err := os.Stat(configPath)
	firstRun := os.IsNotExist(err)

	conf, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to handle config: %v", err)
	}

	if firstRun {
		fmt.Println("First run detected! Created tidy.config.toml.")
		fmt.Println("Please check the target_dir in the config and run the `tidy` again.")
		return
	}

	fmt.Println("Tidying...")
	if err := Tidy(conf); err != nil {
		log.Fatalf("Error during tidy: %v", err)
	}
	fmt.Println("Done!")
}
