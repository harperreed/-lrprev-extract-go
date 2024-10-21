package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/schollz/progressbar/v3"
)

func runGUI() {
	a := app.New()
	w := a.NewWindow("LRPrev Extractor")

	inputDirEntry := widget.NewEntry()
	inputDirEntry.SetPlaceHolder("Path to your lightroom directory (.lrdata)")

	inputFileEntry := widget.NewEntry()
	inputFileEntry.SetPlaceHolder("Path to your file (.lrprev)")

	outputDirEntry := widget.NewEntry()
	outputDirEntry.SetPlaceHolder("Path to output directory")

	lightroomDBEntry := widget.NewEntry()
	lightroomDBEntry.SetPlaceHolder("Path to the lightroom catalog (.lrcat) [optional]")

	includeSizeCheck := widget.NewCheck("Include image size information in the output file name", nil)

	startButton := widget.NewButton("Start", func() {
		inputDir := inputDirEntry.Text
		inputFile := inputFileEntry.Text
		outputDir := outputDirEntry.Text
		lightroomDB := lightroomDBEntry.Text
		includeSize := includeSizeCheck.Checked

		inputPath := inputDir
		if inputFile != "" {
			inputPath = inputFile
		}

		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		fileInfo, err := os.Stat(inputPath)
		if err != nil {
			log.Fatalf("Error accessing input path: %v", err)
		}

		if fileInfo.IsDir() {
			files, err := filepath.Glob(filepath.Join(inputPath, "**/*.lrprev"))
			if err != nil {
				log.Fatalf("Error finding .lrprev files: %v", err)
			}

			bar := progressbar.Default(int64(len(files)))

			for _, file := range files {
				err := processFile(file, outputDir, lightroomDB, includeSize)
				if err != nil {
					fmt.Printf("Error processing file %s: %v\n", file, err)
				}
				if err := bar.Add(1); err != nil {
					fmt.Printf("Error updating progress bar: %v\n", err)
				}
			}
		} else {
			err = processFile(inputPath, outputDir, lightroomDB, includeSize)
			if err != nil {
				log.Fatalf("Error processing file: %v", err)
			}
		}

		fmt.Println("Processing complete!")
	})

	w.SetContent(container.NewVBox(
		inputDirEntry,
		inputFileEntry,
		outputDirEntry,
		lightroomDBEntry,
		includeSizeCheck,
		startButton,
	))

	w.ShowAndRun()
}
