package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cli_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test case 1: Valid path
	err = ValidatePath(tempDir)
	if err != nil {
		t.Errorf("ValidatePath failed for valid path: %v", err)
	}

	// Test case 2: Non-existent path
	nonExistentPath := filepath.Join(tempDir, "non_existent")
	err = ValidatePath(nonExistentPath)
	if err == nil {
		t.Errorf("ValidatePath should have failed for non-existent path")
	}

	// Test case 3: Empty path
	err = ValidatePath("")
	if err == nil {
		t.Errorf("ValidatePath should have failed for empty path")
	}
}

// Note: Testing PromptForInput and PromptForBool would require mocking user input,
// which is beyond the scope of this simple test file. In a real-world scenario,
// you might want to use a mocking library or dependency injection to test these functions.
