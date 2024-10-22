package extractor

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestExtractLargestJPEGFromLRPREV_Success(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "lrprev_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock LRPREV file
	uuid := "12345678-1234-1234-1234-123456789012"
	lrprevPath := filepath.Join(tempDir, "test-"+uuid+".lrprev")
	jpegContent := []byte{0xFF, 0xD8, 0xFF, 0xD9} // Minimal valid JPEG
	err = os.WriteFile(lrprevPath, append([]byte("prefix data"), jpegContent...), 0644)
	assert.NoError(t, err)

	// Run the extraction
	err = ExtractLargestJPEGFromLRPREV(lrprevPath, tempDir, "", false)
	assert.NoError(t, err)

	// Check if the JPEG was extracted correctly - note we're now checking for just the UUID in the output filename
	extractedPath := filepath.Join(tempDir, uuid+".jpg")
	_, err = os.Stat(extractedPath)
	assert.NoError(t, err)

	extractedContent, err := os.ReadFile(extractedPath)
	assert.NoError(t, err)
	assert.Equal(t, jpegContent, extractedContent)
}

func TestExtractLargestJPEGFromLRPREV_InvalidFilePath(t *testing.T) {
	err := ExtractLargestJPEGFromLRPREV("non_existent_file.lrprev", "output", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading file")
}

func TestExtractLargestJPEGFromLRPREV_NoValidJPEG(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "lrprev_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock LRPREV file without a valid JPEG
	uuid := "12345678-1234-1234-1234-123456789012"
	lrprevPath := filepath.Join(tempDir, "test-"+uuid+".lrprev")
	err = os.WriteFile(lrprevPath, []byte("not a valid JPEG"), 0644)
	assert.NoError(t, err)

	// Run the extraction
	err = ExtractLargestJPEGFromLRPREV(lrprevPath, tempDir, "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no valid JPEG found in file")
}

func TestExtractLargestJPEGFromLRPREV_WithDatabase(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "lrprev_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock SQLite database
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := sql.Open("sqlite3", dbPath)
	assert.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE AgLibraryFile (id_global TEXT, folder INTEGER, baseName TEXT);
		CREATE TABLE AgLibraryFolder (id_local INTEGER, rootFolder INTEGER, pathFromRoot TEXT);
		CREATE TABLE AgLibraryRootFolder (id_local INTEGER, absolutePath TEXT);
		INSERT INTO AgLibraryFile (id_global, folder, baseName) VALUES ('12345678-1234-1234-1234-123456789012', 1, 'test');
		INSERT INTO AgLibraryFolder (id_local, rootFolder, pathFromRoot) VALUES (1, 1, 'path/from/root');
		INSERT INTO AgLibraryRootFolder (id_local, absolutePath) VALUES (1, '/absolute/path');
	`)
	assert.NoError(t, err)

	// Create a mock LRPREV file
	uuid := "12345678-1234-1234-1234-123456789012"
	lrprevPath := filepath.Join(tempDir, "test-"+uuid+".lrprev")
	jpegContent := []byte{0xFF, 0xD8, 0xFF, 0xD9} // Minimal valid JPEG
	err = os.WriteFile(lrprevPath, append([]byte("prefix data"), jpegContent...), 0644)
	assert.NoError(t, err)

	// Run the extraction
	err = ExtractLargestJPEGFromLRPREV(lrprevPath, tempDir, dbPath, false)
	assert.NoError(t, err)

	// Check if the JPEG was extracted to the correct path
	extractedPath := filepath.Join(tempDir, "absolute", "path", "path", "from", "root", "test.jpg")
	_, err = os.Stat(extractedPath)
	assert.NoError(t, err)
}

func TestExtractLargestJPEGFromLRPREV_WithoutDatabase(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "lrprev_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock LRPREV file
	uuid := "12345678-1234-1234-1234-123456789012"
	lrprevPath := filepath.Join(tempDir, "test-"+uuid+".lrprev")
	jpegContent := []byte{0xFF, 0xD8, 0xFF, 0xD9} // Minimal valid JPEG
	err = os.WriteFile(lrprevPath, append([]byte("prefix data"), jpegContent...), 0644)
	assert.NoError(t, err)

	// Run the extraction
	err = ExtractLargestJPEGFromLRPREV(lrprevPath, tempDir, "", false)
	assert.NoError(t, err)

	// Check if the JPEG was extracted to the correct path - note we're checking for just the UUID
	extractedPath := filepath.Join(tempDir, uuid+".jpg")
	_, err = os.Stat(extractedPath)
	assert.NoError(t, err)
}

func TestExtractLargestJPEGFromLRPREV_IncludeSize(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "lrprev_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock LRPREV file with a valid JPEG (including JPEG header with dimensions)
	uuid := "12345678-1234-1234-1234-123456789012"
	lrprevPath := filepath.Join(tempDir, "test-"+uuid+".lrprev")
	jpegContent := []byte{
		0xFF, 0xD8, // SOI marker
		0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, // JFIF header
		0xFF, 0xC0, 0x00, 0x11, 0x08, 0x00, 0x10, 0x00, 0x10, 0x03, 0x01, 0x22, 0x00, 0x02, 0x11, 0x01, 0x03, 0x11, 0x01, // SOF marker (16x16 image)
		0xFF, 0xD9, // EOI marker
	}
	err = os.WriteFile(lrprevPath, append([]byte("prefix data"), jpegContent...), 0644)
	assert.NoError(t, err)

	// Run the extraction
	err = ExtractLargestJPEGFromLRPREV(lrprevPath, tempDir, "", true)
	assert.NoError(t, err)

	// Check if the JPEG was extracted with the correct filename (including dimensions)
	extractedPath := filepath.Join(tempDir, uuid+"_16x16.jpg")
	_, err = os.Stat(extractedPath)
	assert.NoError(t, err)
}

func TestExtractLargestJPEGFromLRPREV_ExcludeSize(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "lrprev_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock LRPREV file
	uuid := "12345678-1234-1234-1234-123456789012"
	lrprevPath := filepath.Join(tempDir, "test-"+uuid+".lrprev")
	jpegContent := []byte{0xFF, 0xD8, 0xFF, 0xD9} // Minimal valid JPEG
	err = os.WriteFile(lrprevPath, append([]byte("prefix data"), jpegContent...), 0644)
	assert.NoError(t, err)

	// Run the extraction
	err = ExtractLargestJPEGFromLRPREV(lrprevPath, tempDir, "", false)
	assert.NoError(t, err)

	// Check if the JPEG was extracted with the correct filename (without dimensions)
	extractedPath := filepath.Join(tempDir, uuid+".jpg")
	_, err = os.Stat(extractedPath)
	assert.NoError(t, err)
}
