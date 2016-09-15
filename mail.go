package mail

import (
	"github.com/ishail/m-mail/common"
	"github.com/ishail/m-mail/message"
)

func NewMessage(settings ...message.MessageSetting) *message.Message {
	msg := &message.Message{
		Header:   make(common.Header),
		Charset:  "UTF-8",
		Encoding: common.QuotedPrintable,
	}

	msg.ApplySettings(settings)

	if msg.Encoding == common.Base64 {
		msg.HEncoder = common.BEncoding
	} else {
		msg.HEncoder = common.QEncoding
	}

	return msg
}
