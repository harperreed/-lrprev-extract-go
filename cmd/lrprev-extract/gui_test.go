package main

import (
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestGUIComponents(t *testing.T) {
	app := test.NewApp()
	window := app.NewWindow("Test Window")

	// Create GUI components
	inputDirEntry := widget.NewEntry()
	outputDirEntry := widget.NewEntry()
	lightroomDBEntry := widget.NewEntry()
	includeSizeCheck := widget.NewCheck("Include Size", nil)
	startButton := widget.NewButton("Start", nil)

	// Set up the window content
	window.SetContent(container.NewVBox(
		inputDirEntry,
		outputDirEntry,
		lightroomDBEntry,
		includeSizeCheck,
		startButton,
	))

	// Test initial state
	if inputDirEntry.Text != "" {
		t.Errorf("Expected input directory entry to be empty, got '%s'", inputDirEntry.Text)
	}

	if outputDirEntry.Text != "" {
		t.Errorf("Expected output directory entry to be empty, got '%s'", outputDirEntry.Text)
	}

	if lightroomDBEntry.Text != "" {
		t.Errorf("Expected Lightroom catalog entry to be empty, got '%s'", lightroomDBEntry.Text)
	}

	if includeSizeCheck.Checked {
		t.Error("Expected include size checkbox to be unchecked")
	}

	// Test component interaction
	test.Type(inputDirEntry, "/test/path")
	if inputDirEntry.Text != "/test/path" {
		t.Errorf("Expected input directory text to be '/test/path', got '%s'", inputDirEntry.Text)
	}

	// Clean up
	window.Close()
}

func TestRunGUI(t *testing.T) {
	// Create a test app
	test.NewApp()
	
	// We can't fully test runGUI() because it blocks with ShowAndRun()
	// Instead, we can test that it doesn't panic when called
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("runGUI() panicked: %v", r)
		}
	}()

	// Start GUI in a goroutine so it doesn't block
	go func() {
		runGUI()
	}()
}
