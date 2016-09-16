package sender

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
)

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS {
		advertised := false
		for _, mechanism := range server.Auth {
			if mechanism == "LOGIN" {
				advertised = true
				break
			}
		}
		if !advertised {
			return "", nil, errors.New("m-mail: unencrypted connection")
		}
	}
	if server.Name != a.host {
		return "", nil, errors.New("m-mail: wrong host name")
	}
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	switch {
	case bytes.Equal(fromServer, []byte("Username:")):
		return []byte(a.username), nil
	case bytes.Equal(fromServer, []byte("Password:")):
		return []byte(a.password), nil
	default:
		return nil, fmt.Errorf("m-mail: unexpected server challenge: %s", fromServer)
	}
}
