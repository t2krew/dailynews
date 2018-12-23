package mail

import (
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"sync"

	"github.com/t2krew/daily/output"
)

type EmailError string

func (ee EmailError) Error() string { return string(ee) }

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

func New(email, password, nickname, host string, port int) *mail {
	return &mail{
		smtpHost: host,
		smtpPort: port,
		email:    email,
		password: password,
		nickname: nickname,
	}
}

const contentType = "Content-Type: text/plain; charset=UTF-8"

func (m *mail) Auth() {
	auth := smtp.PlainAuth("", m.email, m.password, m.smtpHost)
	m.auth = &auth
}

func (m *mail) Send(tplname string, receiver []string, content output.Content) (err error) {
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

	message, err := ParseTemplate(fmt.Sprintf("templates/mail/%s.html", tplname), content.Data)
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
	err = smtp.SendMail(addr, *m.auth, m.email, receiver, msgBody)
	if err != nil {
		return
	}

	log.Printf("[task][邮件] done\n")
	return
}
