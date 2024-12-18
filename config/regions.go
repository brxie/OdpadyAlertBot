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
	CalendarFile string `json:"calendar_file"`
	ScheduleTime string `json:"schedule_time"`
}

func LoadRegions() ([]Region, error) {
	// Get the regions file name from the environment variable, default to regionsFile if not set
	fileName := os.Getenv("REGIONS_FILE")
	if fileName == "" {
		fileName = regionsFile
	}

	regionsFile, err := os.Open(fileName)
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
