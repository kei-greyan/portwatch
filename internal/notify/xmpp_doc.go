// Package notify provides notifier implementations for delivering portwatch
// alerts to various external services.
//
// # XMPP Notifier
//
// NewXMPP returns a Notifier that sends alert messages to an XMPP (Jabber)
// recipient via a specified server.
//
// Configuration example (portwatch.yaml):
//
//	  xmpp:
//	    enabled: true
//	    host: "xmpp.example.com"
//	    port: 5222
//	    from: "portwatch@example.com"
//	    password: "s3cr3t"
//	    to: "admin@example.com"
//
// The notifier dials the server for every alert and closes the connection
// immediately after delivery to avoid holding idle sockets.
package notify
