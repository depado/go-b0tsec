package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Check if a folder exists, if not, create it. Returns an error in case of error and nil otherwise.
func CheckAndCreateFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		log.Printf("Could not find %v folder. Creating it.\n", folderPath)
		os.Mkdir(folderPath, 0777)
	} else if err != nil {
		return err
	}
	return nil
}

// Download file and write it disk, using a specific filename
func DownloadNamedFile(url, desiredFilename string) error {
	output, err := os.Create(desiredFilename)
	if err != nil {
		return err
	}
	defer output.Close()
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if _, err := io.Copy(output, response.Body); err != nil {
		return err
	}
	return nil
}

// Download a file and write to disk
func DownloadFile(url string) (string, error) {
	splittedFileName := strings.Split(url, "/")
	fileName := splittedFileName[len(splittedFileName)-1]
	err := DownloadNamedFile(url, fileName)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
