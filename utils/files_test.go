package utils

import (
	"os"
	"testing"
)

func TestCheckAndCreateFolder(t *testing.T) {
	var err error
	testfolder := "testfolder"
	if err = CheckAndCreateFolder(testfolder); err != nil {
		t.Error("Unable to create folder")
	}
	if _, err = os.Stat(testfolder); os.IsNotExist(err) {
		t.Error("Folder was not created")
	}
	if err = CheckAndCreateFolder(testfolder); err != nil {
		t.Error("Error when folder already exists")
	}
	if err = os.Remove(testfolder); err != nil {
		t.Error("Could not remove test folder")
	}
}

func TestHumanReadableSize(t *testing.T) {
	payload := map[int]string{
		1024:  "1.00KB",
		212:   "212B",
		2048:  "2.00KB",
		10000: "9.77KB",
		10e6:  "9.54MB",
		10e10: "93.13GB",
	}
	var ret string
	for key, val := range payload {
		ret = HumanReadableSize(key)
		if ret != val {
			t.Errorf("For %v : Expected %v, got %v", key, val, ret)
		}
	}
}
