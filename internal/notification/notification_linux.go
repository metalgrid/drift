//go:build linux
// +build linux

package notification

import "fmt"

type linuxNotifier struct{}

func (n linuxNotifier) SendNotification() {
	fmt.Println("Linux notification")
}

func newNotifier() Notifier {
	return &linuxNotifier{}
}
