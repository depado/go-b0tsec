package utils

import (
	"fmt"
	"log"
	"os"
)

// CheckAndCreateFolder checks if a folder exists, if not, creates it.
func CheckAndCreateFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		log.Printf("Could not find %v folder. Creating it.\n", folderPath)
		if err := os.Mkdir(folderPath, 0777); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

// HumanReadableSize returns a string with the proper size prefix and with a floating precision of 2
func HumanReadableSize(size int) string {
	if size >= 1024*1024*1024 {
		return fmt.Sprintf("%.2fGB", (float32(size) / 1024.0 / 1024.0 / 1024.0))
	}
	if size >= 1024*1024 {
		return fmt.Sprintf("%.2fMB", (float32(size) / 1024.0 / 1024.0))
	}
	if size >= 1024 {
		return fmt.Sprintf("%.2fKB", (float32(size) / 1024.0))
	}
	return fmt.Sprintf("%dB", size)
}
