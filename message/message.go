/*
	Package message contains Message structure and implements every method on it.
*/
package message

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/ishail/m-mail/common"
)

func NewMessage(subject, body, emailType string, settings ...MessageSetting) *Message {
	if emailType == "html" || emailType == "text/html" {
		emailType = "text/html"
	} else {
		emailType = "text/plain"
	}

	msg := &Message{
		subject:   subject,
		body:      body,
		emailType: emailType,
		header:    make(common.Header),
		charset:   "UTF-8",
		encoding:  common.QuotedPrintable,
	}

	msg.ApplySettings(settings)

	if msg.encoding == common.Base64 {
		msg.hEncoder = common.BEncoding
	} else {
		msg.hEncoder = common.QEncoding
	}

	return msg
}

func (msg *Message) ApplySettings(settings []MessageSetting) {
	for _, setting := range settings {
		setting(msg)
	}
}

// Reset resets the message so it can be reused. The message keeps its previous
// settings so it is in the same state that after a call to NewMessage.
func (msg *Message) Reset() {
	for key := range msg.header {
		delete(msg.header, key)
	}
	msg.parts = nil
	msg.attachments = nil
	msg.embedded = nil
}

func (msg *Message) SetHeader(field string, value ...string) {
	msg.encodeHeader(value)
	msg.header[field] = value
}

func (msg *Message) encodeHeader(values []string) {
	for index, val := range values {
		values[index] = msg.encodeString(val)
	}
}

func (msg *Message) encodeString(value string) string {
	return msg.hEncoder.Encode(msg.charset, value)
}

// SetHeaders sets the message headers.
func (msg *Message) SetHeaders(headers common.Header) {
	for key, val := range headers {
		msg.SetHeader(key, val...)
	}
}

// SetAddressHeader sets an address to the given header field.
func (msg *Message) SetAddressHeader(field, address, name string) {
	msg.header[field] = []string{msg.FormatAddress(address, name)}
}

// FormatAddress formats an address and a name as a valid RFC 5322 address.
func (msg *Message) FormatAddress(address, name string) string {
	if name == "" {
		return address
	}

	enc := msg.encodeString(name)
	if enc == name {
		msg.buff.WriteByte('"')
		for _, character := range name {
			if character == '\\' || character == '"' {
				msg.buff.WriteByte('\\')
			}
			msg.buff.WriteByte(byte(character))
		}
		msg.buff.WriteByte('"')
	} else if common.HasSpecials(name) {
		msg.buff.WriteString(common.BEncoding.Encode(msg.charset, name))
	} else {
		msg.buff.WriteString(enc)
	}

	msg.buff.WriteString(" <")
	msg.buff.WriteString(address)
	msg.buff.WriteByte('>')

	addr := msg.buff.String()
	msg.buff.Reset()
	return addr
}

// SetDateHeader sets a date to the given header field.
func (msg *Message) SetDateHeader(field string, date time.Time) {
	msg.header[field] = []string{common.FormatDate(date)}
}

// GetHeader gets a header field.
func (msg *Message) GetHeader(field string) []string {
	return msg.header[field]
}

//Get From address from Message model
func (msg *Message) GetFrom() (string, error) {
	if from, ok := msg.header["From"]; ok {
		if len(from) > 0 {
			return common.ParseAddress(from[0])
		}
	}
	return "", errors.New("m-mail: invalid message, 'From' field is missing!")
}

//Get list of recipients(To, Cc, Bcc) from Message object
func (msg *Message) GetRecipients() ([]string, error) {
	recipientLength := 0
	addrHeaderList := []string{"To", "Cc", "Bcc"}

	for _, field := range addrHeaderList {
		if addresses, ok := msg.header[field]; ok {
			recipientLength += len(addresses)
		}
	}
	recipients := make([]string, recipientLength)
	index := 0

	for _, field := range addrHeaderList {
		if addresses, ok := msg.header[field]; ok {
			for _, addr := range addresses {
				if addr, err := common.ParseAddress(addr); err != nil {
					return nil, fmt.Errorf(
						"m-mail: Unable to parse address. Address: %s, Error: %v", addr, err)
				} else {
					recipients[index] = addr
					index++
				}
			}
		}
	}

	return recipients, nil
}

//Convert Message object into bytes
func (msg *Message) GetEmailBytes(to string) []byte {
	var msgBytes bytes.Buffer

	msgBytes.WriteString("To: " + to + "\r\n")
	msgBytes.WriteString("Date: " + time.Now().String() + "\r\n")
	msgBytes.WriteString("Subject: " + msg.subject + "\r\n")
	msgBytes.WriteString("Content-Type: multipart/alternative;\r\n")
	msgBytes.WriteString(`    boundary="boundary-type-1234567892-alt"` + "\r\n")
	msgBytes.WriteString("Mime-Version: 1.0\r\n\r\n")
	msgBytes.WriteString("--boundary-type-1234567892-alt\r\n")
	msgBytes.WriteString("Content-Type: " + msg.emailType + `; charset=` + msg.charset + "\r\n")
	msgBytes.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	msgBytes.WriteString(msg.body + "\r\n")

	return msgBytes.Bytes()
}

//Returns headers of message as RFC format
func (msg *Message) getHeadersBytes() []byte {
	var headers bytes.Buffer
	for key, value := range msg.header {
		headers.Write(getHeaderBytes(key, value...))
	}

	return headers.Bytes()
}
