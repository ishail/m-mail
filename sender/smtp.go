package sender

import (
	"crypto/tls"
	"github.com/ishail/m-mail/common"
	"github.com/ishail/m-mail/message"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// NewDialer returns a new SMTP Dialer. The given parameters are used to connect
// to the SMTP server.
func NewDialer(host string, port int, username, password string) *Dialer {
	return &Dialer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		SSL:      port == 465,
	}
}

// Dial dials and authenticates to an SMTP server. The returned SendCloser
// should be closed when done using it.
func (dialer *Dialer) Dial() (SendCloser, error) {
	conn, err := net.DialTimeout("tcp", common.HostPortAddr(dialer.Host, dialer.Port),
		10*time.Second)
	if err != nil {
		return nil, err
	}

	if dialer.SSL {
		conn = tls.Client(conn, dialer.tlsConfig())
	}

	c, err := common.SmtpNewClient(conn, dialer.Host)
	if err != nil {
		return nil, err
	}

	if dialer.LocalName != "" {
		if err := c.Hello(dialer.LocalName); err != nil {
			return nil, err
		}
	}

	if !dialer.SSL {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(dialer.tlsConfig()); err != nil {
				c.Close()
				return nil, err
			}
		}
	}

	if dialer.Auth == nil && dialer.Username != "" {
		if ok, auths := c.Extension("AUTH"); ok {
			if strings.Contains(auths, "CRAM-MD5") {
				dialer.Auth = smtp.CRAMMD5Auth(dialer.Username, dialer.Password)
			} else if strings.Contains(auths, "LOGIN") &&
				!strings.Contains(auths, "PLAIN") {
				dialer.Auth = &loginAuth{
					username: dialer.Username,
					password: dialer.Password,
					host:     dialer.Host,
				}
			} else {
				dialer.Auth = smtp.PlainAuth("", dialer.Username, dialer.Password, dialer.Host)
			}
		}
	}

	if dialer.Auth != nil {
		if err = c.Auth(dialer.Auth); err != nil {
			c.Close()
			return nil, err
		}
	}

	return &smtpSender{*c, dialer}, nil
}

func (dialer *Dialer) tlsConfig() *tls.Config {
	if dialer.TLSConfig == nil {
		return &tls.Config{ServerName: dialer.Host}
	}
	return dialer.TLSConfig
}

func (sender *smtpSender) Send(msg *message.Message) error {
	return nil
}
