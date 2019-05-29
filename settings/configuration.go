package settings

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

//Configuration main configuration params
type Configuration struct {
	Port       string     `yaml:"Port"`
	APIKey     string     `yaml:"ApiKey"`
	SMTPServer SMTPServer `yaml:"SmtpServer"`
	Kafka      Kafka      `yaml:"Kafka"`
}

// SMTPServer configuration for smt sender server
type SMTPServer struct {
	Host           string `yaml:"Host"`
	Port           string `yaml:"Port"`
	SenderEmail    string `yaml:"SenderEmail"`
	SenderPassword string `yaml:"SenderPassword"`
}

// Kafka represents kafka server configuration
type Kafka struct {
	Server string `yaml:"Server"`
	Group  string `yaml:"Group"`
}

// ServerName returns smtp host:port
func (s SMTPServer) ServerName() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

// GetEnv returns env value, if empty, returns fallback value
func GetEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}

//GetConfiguration returns main configuration object
func GetConfiguration() (Configuration, error) {

	configYaml := &Configuration{}
	var err error
	var confFile []byte

	// check if conf-file is .yaml or .json
	confFile, err = ioutil.ReadFile("conf.yaml")
	if err != nil {
		confFile, err = ioutil.ReadFile("conf.json")
		if err != nil {
			log.Println("Could not open conf.yaml or conf.json")
		}
	}

	//if file exists use its variables
	if err == nil {
		err = yaml.Unmarshal(confFile, &configYaml)
		if err != nil {
			return Configuration{}, err
		}
	}

	//Search for OS environment and use conf-file as fallback
	config := Configuration{}
	config.Port = GetEnv("PORT", configYaml.Port)
	config.APIKey = GetEnv("APIKEY", configYaml.APIKey)
	config.SMTPServer.Host = GetEnv("SMTP_HOST", configYaml.SMTPServer.Host)
	config.SMTPServer.Port = GetEnv("SMTP_PORT", configYaml.SMTPServer.Port)
	config.SMTPServer.SenderEmail = GetEnv("SMTP_SENDER_EMAIL", configYaml.SMTPServer.SenderEmail)
	config.SMTPServer.SenderPassword = GetEnv("SMTP_SENDER_PASSWORD", configYaml.SMTPServer.SenderPassword)
	config.Kafka.Server = GetEnv("KAFKA_SERVER", configYaml.Kafka.Server)
	config.Kafka.Group = GetEnv("KAFKA_GROUP", configYaml.Kafka.Group)

	return config, nil
}
