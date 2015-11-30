package usercommand

import (
	"encoding/json"
	"log"

	"github.com/depado/go-b0tsec/database"
)

const bucketName = "usercommand"

// Command represents a single command.
type Command struct {
	Name  string
	Value string
}

// CreateBucket creates the bucket for usercommands.
func CreateBucket() {
	if err := database.BotStorage.CreateBucket(bucketName); err != nil {
		log.Fatalf("While initializing UserCommand plugin : %s", err)
	}
}

// Encode encodes a chain to json.
func (c *Command) Encode() ([]byte, error) {
	enc, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Decode decodes json to Command
func (c *Command) Decode(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}

// Save saves the Data
func (c Command) Save() error {
	return database.BotStorage.Save(bucketName, c.Name, &c)
}

// Delete deletes a command
func (c Command) Delete() error {
	return database.BotStorage.Delete(bucketName, c.Name)
}

// List returns a string array of all the keys in the bucket
func List(list *[]string) error {
	return database.BotStorage.List(bucketName, list)
}
