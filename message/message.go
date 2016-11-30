package message

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/ishail/m-mail/common"
)

func (msg *Message) ApplySettings(settings []MessageSetting) {
	for _, setting := range settings {
		setting(msg)
	}
}

// Reset resets the message so it can be reused. The message keeps its previous
// settings so it is in the same state that after a call to NewMessage.
func (msg *Message) Reset() {
	for key := range msg.Header {
		delete(msg.Header, key)
	}
	msg.Parts = nil
	msg.Attachments = nil
	msg.Embedded = nil
}

func (msg *Message) SetHeader(field string, value ...string) {
	msg.encodeHeader(value)
	msg.Header[field] = value
}

func (msg *Message) encodeHeader(values []string) {
	for index, val := range values {
		values[index] = msg.encodeString(val)
	}
}

func (msg *Message) encodeString(value string) string {
	return msg.HEncoder.Encode(msg.Charset, value)
}

// SetHeaders sets the message headers.
func (msg *Message) SetHeaders(headers common.Header) {
	for key, val := range headers {
		msg.SetHeader(key, val...)
	}
}

// SetAddressHeader sets an address to the given header field.
func (msg *Message) SetAddressHeader(field, address, name string) {
	msg.Header[field] = []string{msg.FormatAddress(address, name)}
}

// FormatAddress formats an address and a name as a valid RFC 5322 address.
func (msg *Message) FormatAddress(address, name string) string {
	if name == "" {
		return address
	}

	enc := msg.encodeString(name)
	if enc == name {
		msg.Buff.WriteByte('"')
		for _, character := range name {
			if character == '\\' || character == '"' {
				msg.Buff.WriteByte('\\')
			}
			msg.Buff.WriteByte(byte(character))
		}
		msg.Buff.WriteByte('"')
	} else if common.HasSpecials(name) {
		msg.Buff.WriteString(common.BEncoding.Encode(msg.Charset, name))
	} else {
		msg.Buff.WriteString(enc)
	}

	msg.Buff.WriteString(" <")
	msg.Buff.WriteString(address)
	msg.Buff.WriteByte('>')

	addr := msg.Buff.String()
	msg.Buff.Reset()
	return addr
}

// SetDateHeader sets a date to the given header field.
func (msg *Message) SetDateHeader(field string, date time.Time) {
	msg.Header[field] = []string{common.FormatDate(date)}
}

// GetHeader gets a header field.
func (msg *Message) GetHeader(field string) []string {
	return msg.Header[field]
}

//Get From address from Message model
func (msg *Message) GetFrom() (string, error) {
	if from, ok := msg.Header["From"]; ok {
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
		if addresses, ok := msg.Header[field]; ok {
			recipientLength += len(addresses)
		}
	}
	recipients := make([]string, recipientLength)
	index := 0

	for _, field := range addrHeaderList {
		if addresses, ok := msg.Header[field]; ok {
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
	msgBytes.WriteString("Subject: " + msg.Subject + "\r\n")
	msgBytes.WriteString("Content-Type: multipart/alternative;\r\n")
	msgBytes.WriteString(`    boundary="boundary-type-1234567892-alt"` + "\r\n")
	msgBytes.WriteString("Mime-Version: 1.0\r\n\r\n")
	msgBytes.WriteString("--boundary-type-1234567892-alt\r\n")
	msgBytes.WriteString("Content-Type: " + msg.Type + `; charset="UTF-8"` + "\r\n")
	msgBytes.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	msgBytes.WriteString(msg.Body + "\r\n")

	return msgBytes.Bytes()
}

//Returns headers of message as RFC format
func (msg *Message) getHeadersBytes() []byte {
	var headers bytes.Buffer
	for key, value := range msg.Header {
		headers.Write(getHeaderBytes(key, value...))
	}

	return headers.Bytes()
}
