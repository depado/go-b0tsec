package utils

import (
	"log"
	"os"
)

// CheckAndCreateFolder checks if a folder exists, if not, creates it.
func CheckAndCreateFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		log.Printf("Could not find %v folder. Creating it.\n", folderPath)
		os.Mkdir(folderPath, 0777)
	} else if err != nil {
		return err
	}
	return nil
}
