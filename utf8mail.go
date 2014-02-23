package utf8mail

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
)

type MailData struct {
	Smtpserver string
	Auth       smtp.Auth
	From       mail.Address
	To         []mail.Address
	Cc         []mail.Address
	Bcc        []mail.Address
	Subject    string
	Body       []byte
}

// 76バイト毎にCRLFを挿入する
func add76crlf(msg string) string {
	var buffer bytes.Buffer
	for k, c := range strings.Split(msg, "") {
		buffer.WriteString(c)
		if k%76 == 75 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String()
}

// UTF8文字列を指定文字数で分割
func utf8Split(utf8string string, length int) []string {
	resultString := []string{}
	var buffer bytes.Buffer
	for k, c := range strings.Split(utf8string, "") {
		buffer.WriteString(c)
		if k%length == length-1 {
			resultString = append(resultString, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 {
		resultString = append(resultString, buffer.String())
	}
	return resultString
}

// サブジェクトをMIMEエンコードする
func encodeSubject(subject string) string {
	var buffer bytes.Buffer
	buffer.WriteString("Subject:")
	for _, line := range utf8Split(subject, 13) {
		buffer.WriteString(" =?utf-8?B?")
		buffer.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		buffer.WriteString("?=\r\n")
	}
	return buffer.String()
}

func (this *MailData) Send() {

	var header bytes.Buffer
	header.WriteString("From: " + this.From.String() + "\r\n")

	header.WriteString("To: ")
	for i, address := range this.To {
		if i != 0 {
			header.WriteString(",")
		}
		header.WriteString(address.String())
	}
	header.WriteString("\r\n")

	header.WriteString("Cc: ")
	for i, cc := range this.Cc {
		if i != 0 {
			header.WriteString(",")
		}
		header.WriteString(cc.String())
	}
	header.WriteString("\r\n")

	header.WriteString(encodeSubject(this.Subject))

	header.WriteString("MIME-Version: 1.0\r\n")
	header.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	header.WriteString("Content-Transfer-Encoding: base64\r\n")

	recpt := []string{}
	for _, address := range this.To {
		recpt = append(recpt, address.Address)
	}
	for _, address := range this.Cc {
		recpt = append(recpt, address.Address)
	}
	for _, address := range this.Bcc {
		recpt = append(recpt, address.Address)
	}

	var message bytes.Buffer
	message = header
	message.WriteString("\r\n")
	message.WriteString(add76crlf(base64.StdEncoding.EncodeToString(this.Body)))

	err := smtp.SendMail(
		this.Smtpserver,
		this.Auth,
		this.From.Address,
		recpt,
		[]byte(message.String()),
	)
	if err != nil {
		log.Fatal(err)
	}
}
