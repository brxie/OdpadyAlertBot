package notification

type NotificationSystem interface {
	SendMessage(chatID, text string) error
}
