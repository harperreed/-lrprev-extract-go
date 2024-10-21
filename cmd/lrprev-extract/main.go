package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"lrprev-extract-go/internal/extractor"
)

func main() {
	inputDir := flag.String("d", "", "Path to your lightroom directory (.lrdata)")
	inputFile := flag.String("f", "", "Path to your file (.lrprev)")
	outputDirectory := flag.String("o", "", "Path to output directory")
	lightroomDB := flag.String("l", "", "Path to the lightroom catalog (.lrcat)")
	includeSize := flag.Bool("include-size", false, "Include image size information in the output file name")
	flag.Parse()

	if *inputDir == "" && *inputFile == "" {
		log.Fatal("Either --input-dir or --input-file must be supplied.")
	}

	if *inputDir != "" && *inputFile != "" {
		log.Fatal("Both --input-dir and --input-file were supplied. Only one is allowed at a time.")
	}

	if *outputDirectory == "" {
		log.Fatal("Output directory must be specified.")
	}

	inputPath := *inputDir
	if *inputFile != "" {
		inputPath = *inputFile
	}

	err := os.MkdirAll(*outputDirectory, os.ModePerm)
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
				return processFile(path, *outputDirectory, *lightroomDB, *includeSize)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error processing directory: %v", err)
		}
	} else {
		err = processFile(inputPath, *outputDirectory, *lightroomDB, *includeSize)
		if err != nil {
			log.Fatalf("Error processing file: %v", err)
		}
	}
}

func processFile(filePath, outputDir, dbPath string, includeSize bool) error {
	fmt.Printf("Processing file: %s\n", filePath)
	return extractor.ExtractLargestJPEGFromLRPREV(filePath, outputDir, dbPath, includeSize)
}
