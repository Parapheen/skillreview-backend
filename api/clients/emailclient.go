package clients

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"net/textproto"
	"os"

	"github.com/jordan-wright/email"
)

type emailClientStruct struct{}

type EmailClientInterface interface {
	Greet(receiver string, nickname string) error
	NewReview(receiver string, nickname string, requestURL string, matchID string) error
	PrepareHTML(fileName string, data interface{}) (*bytes.Buffer, error)
}

var (
	EmailClientStruct EmailClientInterface = &emailClientStruct{}
)

func (cli *emailClientStruct) PrepareHTML(fileName string, data interface{}) (*bytes.Buffer, error) {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return &bytes.Buffer{}, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return &bytes.Buffer{}, err
	}
	return buf, nil
}

func (cli *emailClientStruct) Greet(receiver string, nickname string) error {
	templateData := struct {
		Nickname string
	}{
		Nickname: nickname,
	}
	buf, err := cli.PrepareHTML("emails/greet.html", templateData)
	if err != nil {
		return err
	}
	e := &email.Email{
		To:      []string{receiver},
		From:    fmt.Sprintf("SkillReview <%s>", os.Getenv("EMAIL_USER")),
		Subject: "Welcome to SkillReview",
		HTML:    buf.Bytes(),
		Headers: textproto.MIMEHeader{},
	}

	return e.Send(fmt.Sprintf("%s:587", os.Getenv("EMAIL_SERVER")), smtp.PlainAuth("", os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_SERVER")))
}

func (cli *emailClientStruct) NewReview(receiver string, nickname string, requestURL string, matchID string) error {
	templateData := struct {
		Nickname string
		MatchID  string
		URL      string
	}{
		Nickname: nickname,
		MatchID:  matchID,
		URL:      requestURL,
	}
	buf, err := cli.PrepareHTML("emails/new_review.html", templateData)
	if err != nil {
		return err
	}
	e := &email.Email{
		To:      []string{receiver},
		From:    fmt.Sprintf("SkillReview <%s>", os.Getenv("EMAIL_USER")),
		Subject: "New review for your match",
		HTML:    buf.Bytes(),
		Headers: textproto.MIMEHeader{},
	}

	return e.Send(fmt.Sprintf("%s:587", os.Getenv("EMAIL_SERVER")), smtp.PlainAuth("", os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_SERVER")))
}
