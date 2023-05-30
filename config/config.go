package config

import (
	"encoding/json"
	"io"
	"os"
)

var (
	ListenPort  = "127.0.0.1:7121"
	DBAddress   = "temp" // TODO: implement DB
	AuthAddress = "temp" // TODO: implement Auth
	CorsAllow   = "*"    // TODO: Investigate CORS policy
	Version     = 0.01
	Environment = "Dev"
)

type Config struct {
	ListenPort  string  `json:"listenPort"`
	DBAddress   string  `json:"dbAddress"`
	AuthAddress string  `json:"authAddress"`
	CorsAllow   string  `json:"cors"`
	Version     float32 `json:"version"`
	Environment string  `json:"environment"`
}

// ReadConfig reads the config file for the given environment to set
// default variables for the given environment
func ReadConfig(configFile string) (*Config, error) {
	config := Config{}
	if configFile == "" {
		return nil, nil
	}

	// open configFile
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// unmarshal json to config struct
	json.Unmarshal(b, &config)

	return &config, nil
}
