//go:build linux
// +build linux

package notification

import "fmt"

func SendNotification() {
	fmt.Println("placeholder")
}

type linuxNotifier struct{}

func (n linuxNotifier) SendNotification() {}

func newNotifier() Notifier {
	return &linuxNotifier{}
}
