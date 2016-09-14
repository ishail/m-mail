package message

import (
	"bytes"
	"m-mail/common"
)

func (m *Message) RFCString() ([]byte, error) {
	var response bytes.Buffer
	if _, ok := m.Header["Mime-Version"]; !ok {
		response.WriteString("Mime-Version: 1.0\r\n")
	}

	if _, ok := m.Header["Date"]; !ok {
		response.writeHeader("Date", m.FormatDate(now()))
	}

	return response.Bytes(), nil
}
