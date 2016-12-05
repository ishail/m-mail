package message

import (
	"bytes"

	"github.com/ishail/m-mail/common"
)

// Message represents an email.
type Message struct {
	subject     string
	body        string
	emailType   string
	header      common.Header
	parts       []*common.Part
	attachments []*common.File
	embedded    []*common.File
	charset     string
	encoding    common.Encoding
	hEncoder    common.MimeEncoder
	buff        bytes.Buffer
	trackingUrl string
}

// A MessageSetting can be used as an argument in NewMessage to configure an email.
type MessageSetting func(m *Message)
