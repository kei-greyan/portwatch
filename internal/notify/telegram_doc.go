// Package notify provides notification backends for portwatch alerts.
//
// # Telegram
//
// TelegramNotifier delivers alerts to a Telegram chat using the Bot API.
// Create a bot via @BotFather, obtain the token, and find the numeric chat ID
// of the target chat or channel.
//
// Configuration example (config.yaml):
//
//	notify:
//	  telegram:
//	    enabled: true
//	    token: "123456:ABC-DEF..."
//	    chat_id: "-1001234567890"
//
// Alerts are formatted in Markdown with an emoji prefix indicating severity.
package notify
