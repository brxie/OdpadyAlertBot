package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const regionsFile = "regions.json"

type Region struct {
	ID           string `json:"id"`
	ChatID       string `json:"chat_id"`
	EventsFile   string `json:"events_file"`
	ScheduleTime string `json:"schedule_time"`
}

func LoadRegions() ([]Region, error) {
	regionsFile, err := os.Open(regionsFile)
	if err != nil {
		return nil, fmt.Errorf("can not open regions file: %v", err)
	}

	var r []Region
	jsonParser := json.NewDecoder(regionsFile)
	if err = jsonParser.Decode(&r); err != nil {
		return nil, fmt.Errorf("can not parse regions file: %v", err)
	}

	return r, nil
}
