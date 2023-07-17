package scheduler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/brxie/OdpadyAlertBot/config"
)

const nextTrigCheckBackoff = time.Minute

func Run(region config.Region, events []config.Event, errCh chan error, cb func(e config.Event, scheduledDate time.Time)) {
	regionSchedTime, err := time.Parse("15:04", region.ScheduleTime)
	if err != nil {
		errCh <- fmt.Errorf("unrecoverable scheduler error, can't parse scheduled time, region '%s': %v", region.ID, err)
		return
	}

	for {
		now := time.Now()
		nextSchedCheckTime := time.Date(
			now.Year(), now.Month(), now.Day(),
			regionSchedTime.Hour(), regionSchedTime.Minute(), regionSchedTime.Second(), regionSchedTime.Nanosecond(), time.Local)

		// if trigger time is from past this means it was already sent in this day and
		// we should wait until next day which is nearest possible day of event
		if now.Sub(nextSchedCheckTime) > time.Second*10 {
			nextSchedCheckTime = nextSchedCheckTime.Add(time.Hour * 24)
		}

		log.Printf("next schedule check for region '%s': '%s'", region.ID, nextSchedCheckTime)
		time.Sleep(time.Until(nextSchedCheckTime))

		for _, e := range events {
			eventDate, err := time.Parse("2006-01-02", e.Date)
			if err != nil {
				errCh <- fmt.Errorf("unrecoverable scheduler error, can't event date, region '%s': %v", e.Date, err)
				return
			}

			tommorow := time.Now().AddDate(0, 0, 1)
			if eventDate.Year() == tommorow.Year() &&
				eventDate.Month() == tommorow.Month() &&
				eventDate.Day() == tommorow.Day() {
				log.Printf("sending message for region '%s': %v", region.ID, strings.Join(e.Events, ", "))
				cb(e, nextSchedCheckTime)
				break
			}
		}
		time.Sleep(nextTrigCheckBackoff)
	}
}
