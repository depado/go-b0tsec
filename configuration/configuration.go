package configuration

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Configuration is the main struct that represents a configuration.
type Configuration struct {
	Server      string `yaml:"server"`
	Channel     string `yaml:"channel"`
	BotName     string `yaml:"bot_name"`
	TLS         bool   `yaml:"tls"`
	InsecureTLS bool   `yaml:"insecure_tls"`
	Middlewares []string
	Plugins     []string
	YoutubeKey  string `yaml:"youtube_key"`
}

// Config is the Configuration instance that will be exposed to the other packages.
var Config = new(Configuration)

// Load parses the yml file passed as argument and fills the Config.
func Load(cp string) {
	conf, err := ioutil.ReadFile(cp)
	if err != nil {
		log.Fatalf("Could not read configuration : %v", err)
	}
	err = yaml.Unmarshal(conf, &Config)
	if err != nil {
		log.Fatalf("Error parsing YAML :  %v", err)
	}
}
