package common

import (
	"mime"
	"mime/quotedprintable"
	"net"
	"strings"

	"github.com/ishail/smtp/smtp"
)

var (
	BEncoding = MimeEncoder{mime.BEncoding}
	QEncoding = MimeEncoder{mime.QEncoding}

	LastIndexByte = strings.LastIndexByte
	NewQPWriter   = quotedprintable.NewWriter

	SmtpNewClient = func(conn net.Conn, host string) (*smtp.Client, error) {
		return smtp.NewClient(conn, host)
	}
)
