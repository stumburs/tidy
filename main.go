package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	statsFlag := flag.Bool("stats", false, "Show archive stats")
	searchFlag := flag.String("search", "", "Search for a file by its original name")
	configPath := "tidy.config.toml"

	flag.Parse()

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

	// Flags
	if *statsFlag {
		ShowStats(conf)
		return
	}

	if *searchFlag != "" {
		SearchFiles(conf, *searchFlag)
		return
	}

	fmt.Println("Tidying...")
	if err := Tidy(conf); err != nil {
		log.Fatalf("Error during tidy: %v", err)
	}
	fmt.Println("Done!")
}
