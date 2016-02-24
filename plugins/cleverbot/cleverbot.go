package cleverbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/depado/go-b0tsec/configuration"
)

type CleverBot struct {
	Nick string `json:"nick"`
	User string `json:"user"`
	Key  string `json:"key"`
}

var Clever CleverBot

func (c *CleverBot) Initialize() error {
	c.Nick = "gobot"
	c.User = configuration.Config.CleverBotUser
	c.Key = configuration.Config.CleverBotKey
	if c.Key != "" && c.User != "" {
		return c.Create()
	}
	return nil
}

func (c *CleverBot) Create() error {
	var err error
	var m []byte
	if m, err = json.Marshal(c); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://cleverbot.io/1.0/create", bytes.NewBuffer(m))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var rsp CleverBotCreationResponse
	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Could not create the bot : %s", rsp.Status)
	}
	return err
}

func (c *CleverBot) Query(text string) (string, error) {
	var err error
	var m []byte

	r := CleverBotRequest{c.Nick, c.User, c.Key, text}

	if m, err = json.Marshal(r); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://cleverbot.io/1.0/ask", bytes.NewBuffer(m))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var rsp CleverBotQueryResponse
	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Could not query the bot : %s", rsp.Status)
	}
	return rsp.Response, err
}

type CleverBotRequest struct {
	Nick string `json:"nick"`
	User string `json:"user"`
	Key  string `json:"key"`
	Text string `json:"text"`
}

type CleverBotCreationResponse struct {
	Status string `json:"status"`
	Nick   string `json:"nick"`
}

type CleverBotQueryResponse struct {
	Status   string `json:"status"`
	Response string `json:"response"`
}
