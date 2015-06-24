package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Configuration struct {
	Server  string
	Channel string
	BotName string
}

var Config = new(Configuration)

func LoadConfiguration() {
	conf, err := ioutil.ReadFile("conf/conf.yml")
	if err != nil {
		log.Fatalf("Could not read configuration : %v", err)
	}
	err = yaml.Unmarshal(conf, &Config)
	if err != nil {
		log.Fatalf("Error parsing YAML :  %v", err)
	}
}
