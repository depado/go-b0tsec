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
func InitLoggers() (err error) {
	sc := strings.Replace(configuration.Config.Channel, "#", "", 1)
	err = CheckAndCreateFolder("logs")
	if err != nil {
		return
	}
	lf := path.Join("logs", sc)
	err = CheckAndCreateFolder(lf)
	if err != nil {
		return
	}

	HistoryFile, err := os.OpenFile(path.Join(lf, "history.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	LinkFile, err := os.OpenFile(path.Join(lf, "links.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	LinkLogger = log.New(LinkFile, "", log.Ldate|log.Ltime)
	HistoryLogger = log.New(HistoryFile, "", log.Ldate|log.Ltime)

	return
}
