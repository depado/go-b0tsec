package database

import (
	"time"

	"github.com/boltdb/bolt"
)

// Storage is a type that contains a bolt.DB and a boolean that indicates if the connection is already open or not.
type Storage struct {
	DB     *bolt.DB
	Opened bool
}

// Open opens the database connection and create the file if necessary
func (s *Storage) Open() error {
	var err error
	dbfile := "data.db"
	config := &bolt.Options{Timeout: 1 * time.Second}
	s.DB, err = bolt.Open(dbfile, 0600, config)
	if err != nil {
		return err
	}
	s.Opened = true
	return nil
}

// Close closes the connection (or at least attempts to)
func (s *Storage) Close() error {
	s.Opened = false
	err := s.DB.Close()
	return err
}

// BotStorage is the general storage associated to the bot.
// It should be available to any plugin, middleware or any other part of the program.
var BotStorage = Storage{}
