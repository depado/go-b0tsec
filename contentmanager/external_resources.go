package contentmanager

import (
	"github.com/depado/go-b0tsec/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	configurationFolder = "conf/external/active/"
	backupFolder        = "content/backup/"
	contentFolder       = "content/"
)

type ExternalResource struct {
	UpdateInterval time.Duration
	FriendlyName   string
	Url            string
	FileName       string
	Iterations     int
}

type UnparsedExternalResource struct {
	UpdateInterval string
	FriendlyName   string
	Url            string
	FileName       string
}

type AvailableExternalResource struct {
	Resource ExternalResource
	FullPath string
}

var AvailableExternalResources = make(map[string]AvailableExternalResource)

func LoadAndStartExternalResources() error {
	files, err := ioutil.ReadDir(configurationFolder)
	if err != nil {
		log.Println("Could not read configuration folder :", err)
		return err
	}
	for _, fd := range files {
		log.Println("Loading Configuration :", fd.Name())
		ext, err := LoadExternalResourceConfiguration(configurationFolder + fd.Name())
		if err != nil {
			log.Printf("Error with file %v. It won't be used. %v\n", fd.Name(), err)
			continue
		}
		log.Printf("Starting External Resource Collection for %v (%v)\n", ext.FriendlyName, fd.Name())
		go PeriodicUpdateExternalResource(ext)
	}
	return nil
}

func LoadExternalResourceConfiguration(configPath string) (ExternalResource, error) {
	conf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("Could not read external resource configuration :", err)
		return ExternalResource{}, err
	}
	unparsedExternal := new(UnparsedExternalResource)
	err = yaml.Unmarshal(conf, &unparsedExternal)
	if err != nil {
		log.Println("Error parsing YAML :", err)
		return ExternalResource{}, err
	}
	duration, err := time.ParseDuration(unparsedExternal.UpdateInterval)
	if err != nil {
		log.Println("Error parsing Duration :", err)
		return ExternalResource{}, err
	}
	external := ExternalResource{
		UpdateInterval: duration,
		FriendlyName:   unparsedExternal.FriendlyName,
		Url:            unparsedExternal.Url,
		FileName:       unparsedExternal.FileName,
	}
	return external, nil
}

func CalculateIteration(ext *ExternalResource, backupFolder string) error {
	files, err := ioutil.ReadDir(backupFolder)
	if err != nil {
		return err
	}
	toplevel := 0
	if len(files) > 0 {
		for _, f := range files {
			level, err := strconv.Atoi(f.Name()[len(f.Name())-1:])
			if err != nil {
				log.Printf("Could not calculate Iteration for %v. It will be set to 0.\n", ext.FriendlyName)
				ext.Iterations = 0
				return nil
			}
			if level > toplevel {
				toplevel = level
			}
		}
		ext.Iterations = toplevel + 1
		return nil
	}
	ext.Iterations = toplevel
	return nil
}

func PeriodicUpdateExternalResource(ext ExternalResource) {
	tmpFileName := contentFolder + ext.FileName + ".tmp"
	currentFileName := contentFolder + ext.FileName
	mapName := ext.FileName[0 : len(ext.FileName)-len(filepath.Ext(ext.FileName))]
	specificBackupFolder := backupFolder + mapName + "/"

	if err := utils.CheckAndCreateFolder(contentFolder); err != nil {
		log.Println(err)
		return
	}

	if err := utils.CheckAndCreateFolder(backupFolder); err != nil {
		log.Println(err)
		return
	}

	if err := utils.CheckAndCreateFolder(specificBackupFolder); err != nil {
		log.Println(err)
		return
	}

	if err := CalculateIteration(&ext, specificBackupFolder); err != nil {
		log.Println("Error calculating Iteration :", err)
		return
	}

	if err := utils.DownloadNamedFile(ext.Url, currentFileName); err != nil {
		log.Println("Error dowloading file :", err)
		return
	}

	tickChan := time.NewTicker(ext.UpdateInterval).C

	log.Println("Resource Collection Started for", ext.FriendlyName)
	AvailableExternalResources[mapName] = AvailableExternalResource{ext, currentFileName}
	log.Printf("Added Available Resource %v as %v\n", ext.FriendlyName, mapName)

	for {
		<-tickChan

		if err := utils.DownloadNamedFile(ext.Url, tmpFileName); err != nil {
			log.Println(err)
			continue
		}

		same, err := utils.SameFileCheck(currentFileName, tmpFileName)
		if err != nil {
			log.Println(err)
			continue
		}

		if same {
			if err := os.Remove(tmpFileName); err != nil {
				log.Println(err)
				continue
			}
		} else {
			if err := os.Rename(currentFileName, specificBackupFolder+ext.FileName+"."+strconv.Itoa(ext.Iterations)); err != nil {
				log.Println(err)
				continue
			}
			ext.Iterations += 1
			if err := os.Rename(tmpFileName, currentFileName); err != nil {
				log.Println(err)
				continue
			}
			log.Printf("New content in %v. Downloaded and replaced.\n", ext.FriendlyName)
		}
	}
}
