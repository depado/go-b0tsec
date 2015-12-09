package database

import (
	"encoding/json"
	"reflect"
	"testing"
)

const bucketName = "testing"

// TestStorage represents a single command.
type TestStorage struct {
	Name  string
	Value string
}

// Encode encodes a chain to json.
func (s *TestStorage) Encode() ([]byte, error) {
	enc, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Decode decodes json to TestStorage
func (s *TestStorage) Decode(data []byte) error {
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}
	return nil
}

// TestOpenCloseDB tests the database Open and Close functions
func TestOpenCloseDB(t *testing.T) {
	if BotStorage.Opened {
		if err := BotStorage.Close(); err != nil {
			t.Error("Error while closing the database or `Opened` attribute set to true while it should not be")
		}
	}
	if BotStorage.Opened {
		t.Error("The attribute `Opened` should be set to false and is not")
	}
	if err := BotStorage.Open(); err != nil {
		t.Error("Error while opening the database")
	}
	if !BotStorage.Opened {
		t.Error("The attribute `Opened` should be set to true and is not")
	}
	/*
		// Induces boltdb crash : panic: runtime error: invalid memory address or nil pointer dereference
		if err := BotStorage.Open(); err.Error() != "timeout" {
			t.Error("Error `timeout` should have been reported while opening already opened db instead of : %v", err)
		}
	*/
	if err := BotStorage.Close(); err != nil {
		t.Error("Error while closing the database")
	}
	if BotStorage.Opened {
		t.Error("The attribute `Opened` should be set to false and is not")
	}
}

// TestCreateBucket tests the CreateBucket function
func TestCreateBucket(t *testing.T) {
	if BotStorage.Opened {
		if err := BotStorage.Close(); err != nil {
			t.Error("Could not close the db :", err)
		}
	}

	if err := BotStorage.CreateBucket(bucketName); err.Error() != "db must be opened before creating bucket" {
		t.Errorf("Error not reported correctly, db is not opened : %v", err.Error())
	}
	if err := BotStorage.Open(); err != nil {
		t.Errorf("Cant open db for testing : %v", err)
	}

	if err := BotStorage.CreateBucket(bucketName); err != nil {
		t.Error("Could not create bucket")
	}
}

// TestSave tests the Save function
func TestSave(t *testing.T) {
	s := TestStorage{Name: "test", Value: "value"}

	if err := BotStorage.Save(bucketName, s.Name, &s); err != nil {
		t.Errorf("Error while saving data to bucket : %v", err)
	}
	if err := BotStorage.Close(); err != nil {
		t.Skip("Unable to close the db properly in save test")
	}

	ts := TestStorage{Name: "test1", Value: "value1"}
	if err := BotStorage.Save(bucketName, ts.Name, &ts); err.Error() != "db must be opened before saving" {
		t.Errorf("Error not reported correctly (db must be opened before operation) : %v", err)
	}
}

// TestGet tests the Get function
func TestGet(t *testing.T) {
	s := TestStorage{Name: "test"}
	if BotStorage.Opened {
		if err := BotStorage.Close(); err != nil {
			t.Skip("Unable to close the db properly in Get test")
		}
	}
	if err := BotStorage.Get(bucketName, s.Name, &s); err.Error() != "Database must be opened first." {
		t.Errorf("Error not reported correctly (db must be opened before operation) : %v", err)
	}

	if err := BotStorage.Open(); err != nil {
		t.Skip("Unable to open the db properly in Get test")
	}

	if err := BotStorage.Get(bucketName, s.Name, &s); err != nil {
		t.Errorf("Decode error while testing the Get function")
	}

	if s.Value != "value" {
		t.Errorf("Wrong value got with Get()")
	}
}

// TestDelete tests the deletion of a command
func TestDelete(t *testing.T) {
	if BotStorage.Opened {
		if err := BotStorage.Close(); err != nil {
			t.Skip("Unable to close the db properly while testing the Delete function")
		}
	}
	s := TestStorage{Name: "test"}
	if err := BotStorage.Delete(bucketName, s.Name); err.Error() != "db must be opened before using it" {
		t.Errorf("Error should be reported : db not opened")
	}
	if err := BotStorage.Open(); err != nil {
		t.Skip("Unable to open db in deletion test")
	}
	if err := BotStorage.Delete(bucketName, s.Name); err != nil {
		t.Errorf("Error while deleting key : %v", err)
	}

	if err := BotStorage.Delete("veryImprobableBucketName", s.Name); err != nil {
		t.Errorf("Error reported on non existing bucket : %v", err)
	}
}

// TestList tests the List function
func TestList(t *testing.T) {
	var list []string

	if BotStorage.Opened {
		if err := BotStorage.Close(); err != nil {
			t.Skip("Unable to close the db properly while testing the List function")
		}
	}

	if err := BotStorage.List(bucketName, &list); err.Error() != "Database must be opened first." {
		t.Errorf("Error should be reported (db not opened) instead of : %v", err)
	}

	if err := BotStorage.Open(); err != nil {
		t.Skip("Unable to open the db properly while testing the List function")
	}

	keys := []string{"t1", "t2", "t3", "t4", "t5"}
	BatchSave(keys)

	if err := BotStorage.List(bucketName, &list); err != nil {
		t.Errorf("Error while listing bucket keys")
	}

	if !reflect.DeepEqual(list, keys) {
		t.Errorf("Incorrect list returned")
	}

	if err := BotStorage.List("reallyImprobableBucketName", &list); err != nil {
		t.Errorf("Error reported instead of `nil` with non existing bucket")
	}
}

func BatchSave(keys []string) error {
	var t TestStorage
	for _, k := range keys {
		t = TestStorage{Name: k}
		if err := BotStorage.Save(bucketName, k, &t); err != nil {
			return err
		}
	}
	return nil
}
