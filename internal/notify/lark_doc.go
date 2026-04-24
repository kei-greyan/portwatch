// Package notify provides Notifier implementations for various alerting
// backends.
//
// # Lark / Feishu
//
// NewLark constructs a Notifier that delivers alerts to a Lark (Feishu)
// incoming-webhook URL.  Obtain the URL from the Lark developer console by
// creating a custom bot in the target group chat.
//
// The message is sent as a plain-text card so it renders in both the desktop
// and mobile clients without requiring any additional permissions.
//
// Example:
//
//	n := notify.NewLark("https://open.larksuite.com/open-apis/bot/v2/hook/<token>")
//	n.Send(alert.Alert{...})
package notify
