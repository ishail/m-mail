package mail

import (
	"fmt"
	"github.com/ishail/m-mail/common"
	"github.com/ishail/m-mail/message"
	"github.com/ishail/m-mail/sender"
)

//Creates a Message object with optional MessageSetting
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

// NewPlainDialer returns a new SMTP Dialer. The given parameters are used to
// connect to the SMTP server.
func NewPlainDialer(host string, port int, username, password string) *sender.Dialer {
	return sender.NewDialer(host, port, username, password)
}

// DialAndSend opens a connection to the SMTP server, sends the given emails and
// closes the connection
func DialAndSend(dialer *sender.Dialer, messages ...*message.Message) error {
	sendCloser, err := dialer.Dial()
	if err != nil {
		return err
	}
	defer sendCloser.Close()

	for index, msg := range messages {
		if err = sendCloser.Send(msg); err != nil {
			return fmt.Errorf("m-mail: could not send email %d: %v", index+1, err)
		}
	}

	return nil
}
