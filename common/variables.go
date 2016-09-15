package common

import (
	"mime"
	"mime/quotedprintable"
	"strings"
)

var (
	BEncoding = MimeEncoder{mime.BEncoding}
	QEncoding = MimeEncoder{mime.QEncoding}

	LastIndexByte = strings.LastIndexByte
	NewQPWriter   = quotedprintable.NewWriter
)
