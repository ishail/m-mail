/*
	Package mail exposes functions to create Message object and send them.
*/
package mail

import (
	"fmt"

	"github.com/ishail/m-mail/common"
	"github.com/ishail/m-mail/message"
	"github.com/ishail/m-mail/sender"
	"github.com/ishail/smtp/smtp"
)

// NewMessage creates a Message object with optional MessageSetting.
func NewMessage(subject, body, emailType string, settings ...message.MessageSetting) *message.Message {
	return message.NewMessage(subject, body, emailType, settings...)
}

// NewPlainDialer returns a new SMTP Dialer. The given parameters are used to connect to the SMTP
// server.
func NewPlainDialer(host string, port int, username, password string) *sender.Dialer {
	return sender.NewDialer(host, port, username, password)
}

// DialAndSend opens a connection to the SMTP server, sends the given emails and closes the
// connection.
func DialAndSend(dialer *sender.Dialer, messages ...*message.Message) error {
	sendCloser, err := dialer.Dial()
	if err != nil {
		return err
	}
	defer sendCloser.Close()

	for index, msg := range messages {
		if resp, err := sendCloser.Send(msg); err != nil {
			fmt.Printf("m-mail: could not send email %d: %v", index+1, err)
		} else {
			fmt.Println("response....", resp)
		}
	}

	return nil
}

// Temporary function to send mail. To be scrapped.
func SendMail(auth *smtp.Auth, to, from string, msg *message.Message, host string, port int) (string, error) {
	_, resp, err := smtp.SendMail(
		common.HostPortAddr(host, port),
		*auth,
		from,
		[]string{to},
		msg.GetEmailBytes(to))
	return resp, err
}

// Returns a plain smtp.auth object.
func PlainAuth(username, password, host string) smtp.Auth {
	return smtp.PlainAuth("", username, password, host)
}
