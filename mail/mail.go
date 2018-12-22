package mail

import (
	"fmt"
	"log"
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

type Data struct {
	Date string              `json:"date"`
	List []map[string]string `json:"list"`
}

type Content struct {
	Subject string `json:"subject"`
	Data    *Data  `json:"data"`
	Mime    string `json:"content_type"`
}

const contentType = "Content-Type: text/plain; charset=UTF-8"

func (m *mail) Auth() {
	auth := smtp.PlainAuth("", m.email, m.password, m.smtpHost)
	m.auth = &auth
}

func (m *mail) Send(template string, receiver []string, content Content) (err error) {
	if m.auth == nil {
		m.lock.Lock()
		defer m.lock.Unlock()
		if m.auth == nil {
			m.Auth()
		}
	}

	if content.Subject == "" {
		return ErrEmptySubject
	}

	contType := contentType
	if content.Mime != "" {
		contType = content.Mime
	}

	message, err := ParseTemplate(template, content.Data)
	if err != nil {
		log.Println("parse error ", err)
		return err
	}

	var (
		to      = strings.Join(receiver, ",")
		from    = fmt.Sprintf("%s<%s>", m.nickname, m.email)
		body    = fmt.Sprintf("To: %s \r\nFrom: %s\r\nSubject: %s\r\n%s\r\n\r\n%s", to, from, content.Subject, contType, message)
		msgBody = []byte(body)
		addr    = net.JoinHostPort(m.smtpHost, strconv.Itoa(m.smtpPort))
	)

	//fmt.Println(addr, *m.auth, m.email, receiver, string(msgBody))

	return smtp.SendMail(addr, *m.auth, m.email, receiver, msgBody)
}
