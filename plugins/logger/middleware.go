package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	middlewareName = "logger"
)

var (
	urlRegex, _ = regexp.Compile("^https?:.*$")

	// LinkLogger is used to aggregate links.
	LinkLogger *log.Logger

	// LinkFile is the file associated to the LinkLogger
	LinkFile *os.File

	//HistoryLogger is used to aggregate all the messages.
	HistoryLogger *log.Logger

	// HistoryFile is the file associated to the HistoryLogger
	HistoryFile *os.File
)

// Middleware is the actual logger.Middleware
type Middleware struct {
	Started bool
}

func init() {
	plugins.Middlewares[middlewareName] = new(Middleware)
}

// Get actually operates on the message
func (m *Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if !m.Started {
		return
	}
	HistoryLogger.Println(from + " : " + message)
	for _, field := range strings.Fields(message) {
		if urlRegex.MatchString(field) {
			LinkLogger.Println(from + " : " + field)
		}
	}
}

// Start starts the middleware and returns any occured error, nil otherwise
func (m *Middleware) Start() error {
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		if err := InitLoggers(); err != nil {
			return fmt.Errorf("Error init loggers : %s\n", err)
		}
		m.Started = true
	}
	return nil
}

// Stop returns nil when it didnt encounter any error, the error otherwise
func (m *Middleware) Stop() error {
	var closeErr error
	if err := HistoryFile.Close(); err != nil {
		log.Printf("Error closing history file : %s\n", err)
		closeErr = fmt.Errorf("error occured while closing loggers’ files")
	}
	if err := LinkFile.Close(); err != nil {
		log.Printf("closing links file : %s\n", err)
		closeErr = fmt.Errorf("error occured while closing loggers’ files")
	}

	if closeErr != nil {
		return closeErr
	}
	m.Started = false
	return nil
}

// IsStarted returns the state of the middleware
func (m *Middleware) IsStarted() bool {
	return m.Started
}

// InitLoggers initialize the loggers to use with the logger middleware.
func InitLoggers() error {
	var err error
	sc := strings.Replace(configuration.Config.Channel, "#", "", 1)
	if err = utils.CheckAndCreateFolder("logs"); err != nil {
		return err
	}
	lf := path.Join("logs", sc)
	if err = utils.CheckAndCreateFolder(lf); err != nil {
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
