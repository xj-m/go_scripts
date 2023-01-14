// package main
package main

import (
	"fmt"

	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func main() {
	wechaty.NewWechaty().
		OnScan(func(_ *wechaty.Context, qrCode string, status schemas.ScanStatus, _ string) {
			fmt.Printf("Scan QR Code to login: %s\nhttps://wechaty.github.io/qrcode/%s\n", status, qrCode)
		}).
		OnLogin(func(_ *wechaty.Context, user *user.ContactSelf) {
			fmt.Printf("User %s logined\n", user)
		}).
		OnMessage(func(_ *wechaty.Context, message *user.Message) {
			fmt.Printf("Message: %s\n", message)
		}).DaemonStart()
}
