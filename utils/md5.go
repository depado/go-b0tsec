package utils

import (
	"crypto/md5"
	"io"
	"math"
	"os"
)

// We set the filechunk to 8kb
const filechunk = 8192

// Generate the md5sum of a file.
func GenerateMd5Sum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		file.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return string(hash.Sum(nil)), nil
}

// Calculates if two files have the same checksum (whether they are the same or not)
func SameFileCheck(firstFilename, secondFilename string) (bool, error) {
	firstChecksum, err := GenerateMd5Sum(firstFilename)
	if err != nil {
		return false, err
	}
	secondChecksum, err := GenerateMd5Sum(secondFilename)
	if err != nil {
		return false, err
	}
	return firstChecksum == secondChecksum, nil
}
