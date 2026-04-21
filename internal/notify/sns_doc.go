// Package notify provides notification backends for portwatch alerts.
//
// # SNS Notifier
//
// SNSNotifier publishes alerts to an AWS Simple Notification Service (SNS)
// topic via an HTTP endpoint. It is compatible with real AWS SNS as well as
// local mock servers such as LocalStack.
//
// Example configuration:
//
//	n := notify.NewSNS(
//		"https://sns.us-east-1.amazonaws.com",
//		"arn:aws:sns:us-east-1:123456789012:portwatch-alerts",
//	)
//
// A custom HTTP client (e.g. with AWS request signing middleware) can be
// injected via the WithSNSHTTPClient option.
package notify
