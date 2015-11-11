package utils

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
)

// LinkLogger is used to aggregate links.
var LinkLogger *log.Logger

// LinkFile is the file associated to the LinkLogger
var LinkFile *os.File

//HistoryLogger is used to aggregate all the messages.
var HistoryLogger *log.Logger

// HistoryFile is the file associated to the HistoryLogger
var HistoryFile *os.File

// InitLoggers initialize the loggers to use with the logger middleware.
func InitLoggers() error {
	var err error
	sc := strings.Replace(configuration.Config.Channel, "#", "", 1)
	if err = CheckAndCreateFolder("logs"); err != nil {
		return err
	}
	lf := path.Join("logs", sc)
	if err = CheckAndCreateFolder(lf); err != nil {
		return err
	}

	HistoryFile, err := os.OpenFile(path.Join(lf, "history.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	LinkFile, err := os.OpenFile(path.Join(lf, "links.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	LinkLogger = log.New(LinkFile, "", log.Ldate|log.Ltime)
	HistoryLogger = log.New(HistoryFile, "", log.Ldate|log.Ltime)

	return err
}
