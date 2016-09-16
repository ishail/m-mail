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

// Sender is the interface that wraps the Send method.
// Send sends an email to the given addresses.
type Sender interface {
	Send(from string, to []string, msg []byte) error
}

// SendCloser is the interface that groups the Send and Close methods.
type SendCloser interface {
	Sender
	Close() error
}
