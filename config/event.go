package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Event struct {
	Date   string   `json:"date"`
	Events []string `json:"events"`
}

func LoadEvents(fileName string) ([]Event, error) {
	eventseFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("can not open events file: %v", err)
	}

	var e []Event
	jsonParser := json.NewDecoder(eventseFile)
	if err = jsonParser.Decode(&e); err != nil {
		return nil, fmt.Errorf("can not parse events file: %v", err)
	}

	return e, nil
}
