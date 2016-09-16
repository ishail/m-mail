package common

import (
	"io"
	"mime"
)

type File struct {
	Name     string
	Header   Header
	CopyFunc func(w io.Writer) error
}

//message header
type Header map[string][]string

type Part struct {
	ContentType string
	Copier      func(io.Writer) error
	Encoding    string
}

// Encoding represents a MIME encoding scheme like quoted-printable or base64.
type Encoding string

type MimeEncoder struct {
	mime.WordEncoder
}
