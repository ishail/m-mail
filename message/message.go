package message

import (
	"github.com/ishail/m-mail/common"
	"time"
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

// SetCharset is a message setting to set the charset of the email.
func SetCharset(charset string) MessageSetting {
	return func(msg *Message) {
		msg.Charset = charset
	}
}

// SetEncoding is a message setting to set the encoding of the email.
func SetEncoding(enc common.Encoding) MessageSetting {
	return func(msg *Message) {
		msg.Encoding = enc
	}
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
	msg.Header[field] = []string{msg.FormatDate(date)}
}

// FormatDate formats a date as a valid RFC 5322 date.
func (msg *Message) FormatDate(date time.Time) string {
	return date.Format(time.RFC1123Z)
}

// GetHeader gets a header field.
func (msg *Message) GetHeader(field string) []string {
	return msg.Header[field]
}
