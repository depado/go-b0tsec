package configuration

import (
	"io/ioutil"
	"log"
	"sort"

	"gopkg.in/yaml.v2"
)

// Configuration is the main struct that represents a configuration.
type Configuration struct {
	Server               string   `yaml:"server"`
	Channel              string   `yaml:"channel"`
	BotName              string   `yaml:"bot_name"`
	TLS                  bool     `yaml:"tls"`
	InsecureTLS          bool     `yaml:"insecure_tls"`
	CommandCharacter     string   `yaml:"command_character"`
	Middlewares          []string `yaml:"middlewares"`
	Plugins              []string `yaml:"plugins"`
	GoogleAPIKey         string   `yaml:"google_api_key"`
	YandexTrnslKey       string   `yaml:"yandex_trnsl_key"`
	YandexDictKey        string   `yaml:"yandex_dict_key"`
	Lang                 string   `yaml:"lang"`
	UserCommandCharacter string   `yaml:"user_command_character"`
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
	sort.Strings(Config.Plugins)
	sort.Strings(Config.Middlewares)
}
