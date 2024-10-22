package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rivo/tview"
	"github.com/schollz/progressbar/v3"
	"lrprev-extract-go/internal/cli"
	"lrprev-extract-go/internal/extractor"
)

var (
	app          *tview.Application
	progressBars map[string]*tview.ProgressBar
	logView      *tview.TextView
	pauseChan    chan bool
	resumeChan   chan bool
	cancelChan   chan bool
	isPaused     bool
	mu           sync.Mutex
)

func main() {
	inputDir := flag.String("d", "", "Path to your lightroom directory (.lrdata)")
	inputFile := flag.String("f", "", "Path to your file (.lrprev)")
	outputDirectory := flag.String("o", "", "Path to output directory")
	lightroomDB := flag.String("l", "", "Path to the lightroom catalog (.lrcat)")
	includeSize := flag.Bool("include-size", false, "Include image size information in the output file name")
	help := flag.Bool("help", false, "Show help information")
	flag.Parse()

	if *help {
		printHelp()
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

	app = tview.NewApplication()
	progressBars = make(map[string]*tview.ProgressBar)
	logView = tview.NewTextView().SetDynamicColors(true).SetScrollable(true).SetChangedFunc(func() {
		app.Draw()
	})

	pauseChan = make(chan bool)
	resumeChan = make(chan bool)
	cancelChan = make(chan bool)
	isPaused = false

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(logView, 0, 1, false)

	if fileInfo.IsDir() {
		files, err := filepath.Glob(filepath.Join(inputPath, "**/*.lrprev"))
		if err != nil {
			log.Fatalf("Error finding .lrprev files: %v", err)
		}

		totalFiles := len(files)
		overallProgressBar := tview.NewProgressBar().SetMax(totalFiles)
		progressBars["overall"] = overallProgressBar
		flex.AddItem(overallProgressBar, 1, 0, false)

		for _, file := range files {
			fileProgressBar := tview.NewProgressBar().SetMax(1)
			progressBars[file] = fileProgressBar
			flex.AddItem(fileProgressBar, 1, 0, false)
		}

		go func() {
			for _, file := range files {
				select {
				case <-pauseChan:
					mu.Lock()
					isPaused = true
					mu.Unlock()
					<-resumeChan
					mu.Lock()
					isPaused = false
					mu.Unlock()
				case <-cancelChan:
					return
				default:
					err := processFile(file, *outputDirectory, *lightroomDB, *includeSize)
					if err != nil {
						logMessage(fmt.Sprintf("Error processing file %s: %v\n", file, err))
					}
					progressBars[file].SetProgress(1)
					progressBars["overall"].SetProgress(progressBars["overall"].GetProgress() + 1)
				}
			}
			logMessage("Processing complete!")
		}()
	} else {
		fileProgressBar := tview.NewProgressBar().SetMax(1)
		progressBars[inputPath] = fileProgressBar
		flex.AddItem(fileProgressBar, 1, 0, false)

		go func() {
			err = processFile(inputPath, *outputDirectory, *lightroomDB, *includeSize)
			if err != nil {
				log.Fatalf("Error processing file: %v", err)
			}
			progressBars[inputPath].SetProgress(1)
			logMessage("Processing complete!")
		}()
	}

	flex.AddItem(tview.NewTextView().SetText("Press 'p' to pause, 'r' to resume, 'c' to cancel, 'h' for help"), 1, 0, false)

	app.SetRoot(flex, true).SetFocus(logView).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'p':
			pauseChan <- true
		case 'r':
			resumeChan <- true
		case 'c':
			cancelChan <- true
			app.Stop()
		case 'h':
			showHelpScreen()
		}
		return event
	})

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

func processFile(filePath, outputDir, dbPath string, includeSize bool) error {
	logMessage(fmt.Sprintf("Processing file: %s\n", filePath))
	return extractor.ExtractLargestJPEGFromLRPREV(filePath, outputDir, dbPath, includeSize)
}

func logMessage(message string) {
	mu.Lock()
	defer mu.Unlock()
	logView.Write([]byte(message + "\n"))
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

func showHelpScreen() {
	helpText := `
lrprev-extract-go: Extract JPEG images from Lightroom preview files

Usage:
  lrprev-extract [options]

Options:
  -d string
        Path to your lightroom directory (.lrdata)
  -f string
        Path to your file (.lrprev)
  -o string
        Path to output directory
  -l string
        Path to the lightroom catalog (.lrcat)
  -include-size
        Include image size information in the output file name
  -help
        Show help information

Examples:
  lrprev-extract -d /path/to/lightroom/directory -o /path/to/output
  lrprev-extract -f /path/to/file.lrprev -o /path/to/output -l /path/to/catalog.lrcat

Press any key to return...
`
	helpView := tview.NewTextView().SetText(helpText).SetDynamicColors(true).SetScrollable(true).SetChangedFunc(func() {
		app.Draw()
	})

	helpView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		app.SetRoot(flex, true).SetFocus(logView)
		return event
	})

	app.SetRoot(helpView, true).SetFocus(helpView)
}
