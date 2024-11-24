//go:build windows
// +build windows

package notification

import "fmt"

type windowsNotifier struct{}

func (n *windowsNotifier) SendNotification() {
	fmt.Printf("Windows notification")
}

func newNotifier() Notifier {
	return &windowsNotifier{}
}
