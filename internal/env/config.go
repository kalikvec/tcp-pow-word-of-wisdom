package env

import (
	"encoding/json"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerHost            string `envconfig:"SERVER_HOST"`
	ServerPort            int    `envconfig:"SERVER_PORT"`
	HashcashZerosCount    int
	ChallengeTimeout      int64
	HashcashMaxIterations int
}

func LoadConfig(path string) (*Config, error) {
	config := Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return &config, err
	}
	err = envconfig.Process("", &config)
	return &config, err
}
