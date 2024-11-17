package notification

type Notifier interface {
	SendNotification()
}

func NewNotifier() Notifier {
	return newNotifier()
}
