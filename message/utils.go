package message

import (
	"bytes"

	"github.com/ishail/m-mail/common"
)

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

//Returns a single header of message as RFC format
func getHeaderBytes(key string, value ...string) []byte {
	var header bytes.Buffer
	header.WriteString(key)
	if len(value) == 0 {
		header.WriteString(":\r\n")
		return header.Bytes()
	}

	header.WriteString(": ")

	// Max header line length is 78 characters in RFC 5322 and 76 characters
	// in RFC 2047. So for the sake of simplicity we use the 76 characters limit.
	// charsLeft := 76 - len(key) - len(": ")

	for index, val := range value {
		// // If the line is already too long, insert a newline right away.
		// if charsLeft < 1 {
		// 	if index == 0 {
		// 		header.WriteString("\r\n ")
		// 	} else {
		// 		header.WriteString(",\r\n ")
		// 	}
		// 	charsLeft = 75
		// } else if index != 0 {
		// 	header.WriteString(", ")
		// 	charsLeft -= 2
		// }
		if index != 0 {
			header.WriteString(", ")
		}
		header.WriteString(val)
	}
	header.WriteString("\r\n")

	return header.Bytes()
}
