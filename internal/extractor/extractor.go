package extractor

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"

	"lrprev-extract-go/internal/database"
	"lrprev-extract-go/internal/utils"
)

func ExtractLargestJPEGFromLRPREV(filePath, outputDir, dbPath string, includeSize bool) error {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	uuid, err := utils.ExtractUUIDFromFilename(filePath)
	if err != nil {
		return err
	}

	jpegStart := bytes.LastIndex(fileContents, []byte{0xFF, 0xD8})
	jpegEnd := bytes.LastIndex(fileContents, []byte{0xFF, 0xD9})

	if jpegStart == -1 || jpegEnd == -1 || jpegEnd <= jpegStart {
		return fmt.Errorf("no valid JPEG found in file")
	}

	jpegContents := fileContents[jpegStart : jpegEnd+2]

	var finalOutputDir string
	var baseName string

	if dbPath != "" {
		originalFilePath, origBaseName, err := database.GetOriginalFilePath(dbPath, uuid)
		if err != nil {
			fmt.Printf("Error getting original file path: %v\n", err)
			finalOutputDir = filepath.Join(outputDir, "_path_not_found")
			baseName = uuid
		} else {
			finalOutputDir = filepath.Join(outputDir, originalFilePath)
			baseName = origBaseName
		}
	} else {
		finalOutputDir = outputDir
		baseName = uuid
	}

	err = os.MkdirAll(finalOutputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	newFilename := fmt.Sprintf("%s.jpg", baseName)

	if includeSize {
		config, err := jpeg.DecodeConfig(bytes.NewReader(jpegContents))
		if err != nil {
			return fmt.Errorf("error decoding JPEG dimensions: %v", err)
		}
		newFilename = fmt.Sprintf("%s_%dx%d.jpg", baseName, config.Width, config.Height)
	}

	jpegPath := filepath.Join(finalOutputDir, newFilename)

	err = ioutil.WriteFile(jpegPath, jpegContents, 0644)
	if err != nil {
		return fmt.Errorf("error writing JPEG file: %v", err)
	}

	fmt.Printf("JPEG image extracted and saved to %s\n", jpegPath)
	return nil
}