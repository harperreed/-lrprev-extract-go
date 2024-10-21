package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"lrprev-extract-go/internal/config"
	"lrprev-extract-go/internal/extractor"
)

func main() {
	configFile := flag.String("c", "", "Path to config file")
	inputDir := flag.String("d", "", "Path to your lightroom directory (.lrdata)")
	inputFile := flag.String("f", "", "Path to your file (.lrprev)")
	outputDirectory := flag.String("o", "", "Path to output directory")
	lightroomDB := flag.String("l", "", "Path to the lightroom catalog (.lrcat)")
	includeSize := flag.Bool("include-size", false, "Include image size information in the output file name")
	flag.Parse()

	var cfg *config.Config
	var err error

	if *configFile != "" {
		cfg, err = config.LoadConfig(*configFile)
		if err != nil {
			log.Fatalf("Error loading config file: %v", err)
		}
	} else {
		cfg = &config.Config{}
	}

	// Merge command-line arguments with config file settings
	if *inputDir != "" {
		cfg.InputDir = *inputDir
	}
	if *inputFile != "" {
		cfg.InputFile = *inputFile
	}
	if *outputDirectory != "" {
		cfg.OutputDirectory = *outputDirectory
	}
	if *lightroomDB != "" {
		cfg.LightroomDB = *lightroomDB
	}
	if *includeSize {
		cfg.IncludeSize = *includeSize
	}

	if cfg.InputDir == "" && cfg.InputFile == "" {
		log.Fatal("Either input directory or input file must be supplied.")
	}

	if cfg.InputDir != "" && cfg.InputFile != "" {
		log.Fatal("Both input directory and input file were supplied. Only one is allowed at a time.")
	}

	if cfg.OutputDirectory == "" {
		log.Fatal("Output directory must be specified.")
	}

	inputPath := cfg.InputDir
	if cfg.InputFile != "" {
		inputPath = cfg.InputFile
	}

	err = os.MkdirAll(cfg.OutputDirectory, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Fatalf("Error accessing input path: %v", err)
	}

	if fileInfo.IsDir() {
		err = filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".lrprev" {
				return processFile(path, cfg.OutputDirectory, cfg.LightroomDB, cfg.IncludeSize)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error processing directory: %v", err)
		}
	} else {
		err = processFile(inputPath, cfg.OutputDirectory, cfg.LightroomDB, cfg.IncludeSize)
		if err != nil {
			log.Fatalf("Error processing file: %v", err)
		}
	}
}

func processFile(filePath, outputDir, dbPath string, includeSize bool) error {
	fmt.Printf("Processing file: %s\n", filePath)
	return extractor.ExtractLargestJPEGFromLRPREV(filePath, outputDir, dbPath, includeSize)
}
