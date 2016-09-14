package common

import (
	"bytes"
	"io"
	"mime"
)

type File struct {
	Name     string
	Header   map[string][]string
	CopyFunc func(w io.Writer) error
}

type Header map[string][]string

// Message represents an email.
type Message struct {
	Header      Header
	Parts       []*Part
	Attachments []*File
	Embedded    []*File
	Charset     string
	Encoding    string
	HEncoder    mime.WordEncoder
	Buff        bytes.Buffer
}

type Part struct {
	ContentType string
	Copier      func(io.Writer) error
	Encoding    string
}
