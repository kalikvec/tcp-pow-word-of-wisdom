package env

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerHost            string
	ServerPort            int
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
	return &config, err
}
