package markovchains

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/depado/go-b0tsec/database"
)

const chainBucketName = "markov"

// Save saves a chain to the database
func (c *Chain) Save() error {
	db := database.BotStorage.DB
	if !database.BotStorage.Opened {
		return fmt.Errorf("db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		mBucket, err := tx.CreateBucketIfNotExists([]byte(chainBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		enc, err := c.Encode()
		if err != nil {
			return fmt.Errorf("Could not encode Chain : %s", err)
		}
		err = mBucket.Put([]byte(c.Key), enc)
		return err
	})
	return err
}

// GetChain gets the markov chain associated to the nick
func GetChain(key string) (*Chain, error) {
	db := database.BotStorage.DB
	if !database.BotStorage.Opened {
		return nil, fmt.Errorf("db must be opened before saving")
	}
	var c *Chain
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(chainBucketName))
		k := []byte(key)
		c, err = Decode(b.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Encode encodes a chain to json.
func (c *Chain) Encode() ([]byte, error) {
	enc, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Decode decodes json to Chain
func Decode(data []byte) (*Chain, error) {
	var c *Chain
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// InitBucketIfNotExists creates the bucket if it doesn't already exists
func InitBucketIfNotExists() error {
	err := database.BotStorage.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(chainBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return err
	})
	return err
}
