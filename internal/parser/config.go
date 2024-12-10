package parser

import (
	"encoding/json"
	"log"
	"os"

	"github.com/KirantiLoh/ssh-automator/internal/model"
)

func ParseConfigFile(fileName string) (*model.Config, error) {
	log.Println("Reading from config file...")
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var res model.Config
	if err := json.NewDecoder(file).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
