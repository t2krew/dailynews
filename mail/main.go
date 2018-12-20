package mail

import (
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
)

type EmailError string

func (ee EmailError) Error() string { return string(ee) }

const ErrEmptyBody = EmailError("email message is empty")
const ErrEmptySubject = EmailError("email subject is empty")

type mail struct {
	email    string
	password string
	nickname string
	smtpPort int
	smtpHost string
	auth     *smtp.Auth
	lock     sync.Mutex
}

func NewMail(email, password, nickname, host string, port int) *mail {
	return &mail{
		smtpHost: host,
		smtpPort: port,
		email:    email,
		password: password,
		nickname: nickname,
	}
}

type Content struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

const contentType = "Content-Type: text/plain; charset=UTF-8"

func (e *mail) Auth() {
	auth := smtp.PlainAuth("", e.email, e.password, e.smtpHost)
	e.auth = &auth
}

func (e *mail) Send(receiver []string, content Content) (err error) {
	if e.auth == nil {
		e.lock.Lock()
		defer e.lock.Unlock()
		if e.auth == nil {
			e.Auth()
		}
	}

	if content.Subject == "" {
		return ErrEmptySubject
	}

	if content.Message == "" {
		return ErrEmptyBody
	}

	var (
		to      = strings.Join(receiver, ",")
		from    = fmt.Sprintf("%s<%s>", e.nickname, e.email)
		body    = fmt.Sprintf("To: %s \r\nFrom: %s\r\nSubject: %s\r\n%s\r\n\r\n%s", to, from, content.Subject, contentType, content.Message)
		msgBody = []byte(body)
		addr    = net.JoinHostPort(e.smtpHost, strconv.Itoa(e.smtpPort))
	)

	fmt.Println(addr, *e.auth, e.email, receiver, string(msgBody))

	return smtp.SendMail(addr, *e.auth, e.email, receiver, msgBody)
}
