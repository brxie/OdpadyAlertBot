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

const notifiFmtStrPL = "♻️ jutro (%s) odbiór odpadów 🚛:\n"

func SchedulerCallback(client *telegram.Telegram, region config.Region) func(e config.Event, scheduledDate time.Time) {

	return func(e config.Event, scheduledDate time.Time) {
		var sb strings.Builder
		weekDay := getWeekdayPL(int(scheduledDate.Weekday()))
		sb.WriteString(fmt.Sprintf(notifiFmtStrPL, weekDay))
		for _, evtName := range e.Events {
			sb.WriteString(fmt.Sprintf("👉 %s\n", evtName))
		}

		if err := client.SendMessage(region.ChatID, sb.String()); err != nil {
			var reqErr *telegram.RequestError
			switch {
			case errors.As(err, &reqErr):
				log.Printf("sending message response error: %s, body: %s", reqErr.Error(), reqErr.ResponseBody)
			default:
				log.Printf("sending message error %v", err)
			}
		}
	}
}

func getWeekdayPL(dayNumber int) string {
	weekdays := []string{
		"Poniedziałek",
		"Wtorek",
		"Środa",
		"Czwartek",
		"Piątek",
		"Sobota",
		"Niedziela"}
	return weekdays[dayNumber]
}
