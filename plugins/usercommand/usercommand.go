package usercommand

import (
	"encoding/json"
	"log"

	"github.com/depado/go-b0tsec/database"
)

const bucketName = "usercommand"

// UserCommand represents a single command.
type UserCommand struct {
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
func (uc *UserCommand) Encode() ([]byte, error) {
	enc, err := json.Marshal(uc)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Decode decodes json to Command
func (uc *UserCommand) Decode(data []byte) error {
	if err := json.Unmarshal(data, uc); err != nil {
		return err
	}
	return nil
}

// Save saves the Data
func (uc UserCommand) Save() error {
	return database.BotStorage.Save(bucketName, uc.Name, &uc)
}

// Delete deletes a command
func (uc UserCommand) Delete() error {
	return database.BotStorage.Delete(bucketName, uc.Name)
}

// List returns a string array of all the keys in the bucket
func List(list *[]string) error {
	return database.BotStorage.List(bucketName, list)
}
