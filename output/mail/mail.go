package mail

import (
	"fmt"
	"log"
	"sync"

	"github.com/t2krew/dailynews/output"
	"gopkg.in/gomail.v2"
)

type EmailError string

func (ee EmailError) Error() string { return string(ee) }

const OutputEmail = "email"
const ErrEmptySubject = EmailError("email subject is empty")

type mail struct {
	output.Output
	email    string
	password string
	nickname string
	smtpPort int
	smtpHost string
	dailer   *gomail.Dialer
	lock     sync.Mutex
}

func New(email, password, nickname, host string, port int) *mail {
	dailer := gomail.NewDialer(host, port, email, password)

	m := &mail{
		email:    email,
		password: password,
		smtpHost: host,
		smtpPort: port,
		nickname: nickname,
		dailer:   dailer,
	}
	m.OutputName = OutputEmail
	return m
}

const contentType = "text/html"

func (m *mail) Name() string {
	return m.OutputName
}

func (m *mail) Send(tplname string, receiver []string, content output.Content) (err error) {
	if content.Subject == "" {
		return ErrEmptySubject
	}

	message, err := ParseTemplate(fmt.Sprintf("templates/mail/%s.html", tplname), content.Data)
	if err != nil {
		log.Println("parse error ", err)
		return err
	}

	email := gomail.NewMessage()
	email.SetHeader("From", m.email, m.nickname)
	email.SetHeader("To", receiver...)
	email.SetHeader("Subject", content.Subject)
	email.SetBody(contentType, message)

	// Send the email to Bob, Cora and Dan.
	if err := m.dailer.DialAndSend(email); err != nil {
		panic(err)
	}

	log.Printf("[task][邮件] done\n")
	return
}
