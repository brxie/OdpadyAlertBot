package scheduler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/brxie/OdpadyAlertBot/config"
)

func Run(region config.Region, events []config.Event, errCh chan error, cb func(e config.Event, scheduledDate time.Time)) {
	log.Printf("started scheduler for region '%s' with schedule time '%s'", region.ID, region.ScheduleTime)
	for {
		t, err := time.Parse("15:04", region.ScheduleTime)
		if err != nil {
			errCh <- fmt.Errorf("unrecoverable scheduler error, can't parse scheduled time, region '%s': %v", region.ID, err)
			return
		}

		now := time.Now()
		todaySchedTime := time.Date(
			now.Year(), now.Month(), now.Day(),
			t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)

		// if trigger time is from past and older than 1m do not schedule to avoid sending
		// notification twice after restart
		if now.Sub(todaySchedTime) > time.Minute {
			time.Sleep(time.Minute)
			continue
		}
		log.Printf("next schedule trigget check for region '%s': '%s'", region.ID, todaySchedTime)
		time.Sleep(time.Until(todaySchedTime))

		for _, e := range events {
			eventDate, err := time.Parse("2006-01-02", e.Date)
			if err != nil {
				errCh <- fmt.Errorf("unrecoverable scheduler error, can't event date, region '%s': %v", e.Date, err)
				return
			}

			tommorow := now.AddDate(0, 0, 1)
			if eventDate.Year() == tommorow.Year() &&
				eventDate.Month() == tommorow.Month() &&
				eventDate.Day() == tommorow.Day() {
				log.Printf("sending message for region '%s': %v", region.ID, strings.Join(e.Events, ", "))
				cb(e, todaySchedTime)
			}
		}
		time.Sleep(time.Minute)
	}
}
