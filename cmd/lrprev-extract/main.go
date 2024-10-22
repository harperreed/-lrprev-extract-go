package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"lrprev-extract-go/internal/cli"
	"lrprev-extract-go/internal/extractor"


	"github.com/rivo/tview"

)

func main() {
	inputDir := flag.String("d", "", "Path to your lightroom directory (.lrdata)")
	inputFile := flag.String("f", "", "Path to your file (.lrprev)")
	outputDirectory := flag.String("o", "", "Path to output directory")
	lightroomDB := flag.String("l", "", "Path to the lightroom catalog (.lrcat)")
	includeSize := flag.Bool("include-size", false, "Include image size information in the output file name")
	help := flag.Bool("help", false, "Show help information")
	gui := flag.Bool("gui", false, "Launch the GUI")
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *gui || (flag.NFlag() == 0 && flag.NArg() == 0) {
		runGUI()
		return
	}

	if *inputDir == "" && *inputFile == "" {
		*inputDir = cli.PromptForInput("Enter the path to your lightroom directory (.lrdata) or file (.lrprev): ")
	}

	if *outputDirectory == "" {
		*outputDirectory = cli.PromptForInput("Enter the path to the output directory: ")
	}

	if *lightroomDB == "" {
		*lightroomDB = cli.PromptForInput("Enter the path to the lightroom catalog (.lrcat) [optional]: ")
	}

	if !*includeSize {
		*includeSize = cli.PromptForBool("Include image size information in the output file name? (y/n): ")
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

	app := tview.NewApplication()
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a text view for logs
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	// Create a gauge for progress
	gauge := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	// Add items to the flex container
	flex.AddItem(gauge, 1, 0, false)
	flex.AddItem(logView, 0, 1, false)

	go func() {
		if fileInfo.IsDir() {
			files, err := filepath.Glob(filepath.Join(inputPath, "**/*.lrprev"))
			if err != nil {
				fmt.Fprintf(logView, "[red]Error finding .lrprev files: %v\n", err)
				app.Stop()
				return
			}

			totalFiles := len(files)
			for i, file := range files {
				progress := int(float64(i+1) / float64(totalFiles) * 100)
				gauge.Clear()
				fmt.Fprintf(gauge, "[yellow]Progress: [white]%d%%", progress)
				app.Draw()

				err := processFile(file, *outputDirectory, *lightroomDB, *includeSize, logView)
				if err != nil {
					fmt.Fprintf(logView, "[red]Error processing file %s: %v\n", file, err)
				}
			}
		} else {
			gauge.Clear()
			fmt.Fprintf(gauge, "[yellow]Progress: [white]0%%")
			err = processFile(inputPath, *outputDirectory, *lightroomDB, *includeSize, logView)
			if err != nil {
				fmt.Fprintf(logView, "[red]Error processing file: %v\n", err)
			}
			fmt.Fprintf(gauge, "[yellow]Progress: [white]100%%")
		}

		fmt.Fprintln(logView, "[green]Processing complete!")
		app.Stop()
	}()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

func processFile(filePath, outputDir, dbPath string, includeSize bool, logView *tview.TextView) error {
	fmt.Fprintf(logView, "Processing file: %s\n", filePath)
	return extractor.ExtractLargestJPEGFromLRPREV(filePath, outputDir, dbPath, includeSize)
}

func printHelp() {
	fmt.Println("lrprev-extract-go: Extract JPEG images from Lightroom preview files")
	fmt.Println("\nUsage:")
	fmt.Println("  lrprev-extract [options]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  lrprev-extract -d /path/to/lightroom/directory -o /path/to/output")
	fmt.Println("  lrprev-extract -f /path/to/file.lrprev -o /path/to/output -l /path/to/catalog.lrcat")
}
