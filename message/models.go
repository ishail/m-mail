package message

import (
	"bytes"
	"github.com/ishail/m-mail/common"
)

// Message represents an email.
type Message struct {
	Subject     string
	Body        string
	Type        string
	Header      common.Header
	Parts       []*common.Part
	Attachments []*common.File
	Embedded    []*common.File
	Charset     string
	Encoding    common.Encoding
	HEncoder    common.MimeEncoder
	Buff        bytes.Buffer
}

// A MessageSetting can be used as an argument in NewMessage to configure an
// email.
type MessageSetting func(m *Message)
