package main

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestRunGUI(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	go runGUI()

	window := app.NewWindow("Test Window")
	window.ShowAndRun()

	if window.Title() != "Lightroom Preview Extractor" {
		t.Errorf("Expected window title to be 'Lightroom Preview Extractor', but got '%s'", window.Title())
	}

	inputDirEntry := test.WidgetRenderer(window.Content().(*fyne.Container).Objects[1]).(*widget.Entry)
	if inputDirEntry.Text != "" {
		t.Errorf("Expected input directory entry to be empty, but got '%s'", inputDirEntry.Text)
	}

	outputDirEntry := test.WidgetRenderer(window.Content().(*fyne.Container).Objects[3]).(*widget.Entry)
	if outputDirEntry.Text != "" {
		t.Errorf("Expected output directory entry to be empty, but got '%s'", outputDirEntry.Text)
	}

	lightroomDBEntry := test.WidgetRenderer(window.Content().(*fyne.Container).Objects[5]).(*widget.Entry)
	if lightroomDBEntry.Text != "" {
		t.Errorf("Expected Lightroom catalog entry to be empty, but got '%s'", lightroomDBEntry.Text)
	}

	includeSizeCheck := test.WidgetRenderer(window.Content().(*fyne.Container).Objects[7]).(*widget.Check)
	if includeSizeCheck.Checked {
		t.Errorf("Expected include size checkbox to be unchecked, but it was checked")
	}
}
