package calendar_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/brxie/OdpadyAlertBot/config"
)

const calendarDir = "calendar"

func TestNoDuplicatedDays(t *testing.T) {
	err := filepath.Walk(calendarDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if err := assertNoDuplicateKey(path); err != nil {
				t.Fatalf("failed date duplicate validation of file '%s': %v", info.Name(), err)
			}

			if err := assertTimeParsed(path); err != nil {
				t.Fatalf("failed time parse validation of file '%s': %v", info.Name(), err)
			}
			return nil
		})
	if err != nil {
		t.Fatal(err)
	}
}

func assertNoDuplicateKey(fullFileName string) error {
	plan, err := ioutil.ReadFile(path.Join(fullFileName))
	if err != nil {
		return err
	}

	var events []config.Event
	if err := json.Unmarshal(plan, &events); err != nil {
		return err
	}

	var dates []string
	for _, e := range events {
		if contains(dates, e.Date) {
			return fmt.Errorf("found duplicated date '%s'", e.Date)
		}
		dates = append(dates, e.Date)
	}

	return nil
}

func assertTimeParsed(fullFileName string) error {
	plan, err := ioutil.ReadFile(path.Join(fullFileName))
	if err != nil {
		return err
	}

	var events []config.Event
	if err := json.Unmarshal(plan, &events); err != nil {
		return err
	}

	for _, e := range events {
		if _, err := time.Parse("2006-01-02", e.Date); err != nil {
			return fmt.Errorf("failed to parse date %s, %v", e.Date, err)
		}
	}

	return nil
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
