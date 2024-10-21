package utils

import (
	"fmt"
	"regexp"
)

func ExtractUUIDFromFilename(filename string) (string, error) {
	uuidPattern := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	match := uuidPattern.FindString(filename)
	if match == "" {
		return "", fmt.Errorf("UUID could not be extracted from the filename: %s", filename)
	}
	return match, nil
}
