package notification

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/brxie/OdpadyAlertBot/config"
	"github.com/brxie/OdpadyAlertBot/pkg/client/telegram"
)

const notifiFmtStrPL = "🔔 jutro (%s) odbiór odpadów:\n\n"

func SchedulerCallback(client NotificationSystem, region config.Region) func(e config.Event, scheduledDate time.Time) {

	return func(e config.Event, scheduledDate time.Time) {
		var sb strings.Builder
		weekDay := getWeekdayPL(int(scheduledDate.Weekday()))
		sb.WriteString(fmt.Sprintf(notifiFmtStrPL, weekDay))
		for _, evtName := range e.Events {
			sb.WriteString(fmt.Sprintf("➡️ %s\n", evtName))
		}

		if err := client.SendMessage(region.ChatID, sb.String()); err != nil {
			var reqErr *telegram.RequestError
			switch {
			case errors.As(err, &reqErr):
				log.Printf("sending message response error: %s, region: '%s'. body: %s", reqErr.Error(), region.ID, reqErr.ResponseBody)
			default:
				log.Printf("sending message error: %v, region: '%s'", err, region.ID)
			}
		}
	}
}

func getWeekdayPL(dayNumber int) string {
	weekdays := []string{
		"poniedziałek",
		"wtorek",
		"środa",
		"czwartek",
		"piątek",
		"sobota",
		"niedziela"}
	return weekdays[dayNumber]
}
