utf8mail
========


```go:example.go
package main

import (
	"bytes"
	"github.com/masahide/utf8mail"
	"net/mail"
)

func main() {

	var body bytes.Buffer
	body.WriteString("本文本文本文本文本文本文本文本文本文本文本文本文本文本文本文本文本文本文\n")
	body.WriteString("本文本文本文本文本文本文本文本文本文本文本文本文本文本文本文本文本文本文\n")

	mail := utf8mail.MailData{
		"localhost:25",
		nil, //auth
		mail.Address{"ほげほげ", "masahide.y@gmail.com"},     //from
		[]mail.Address{{"ふがふが", "masahide.y@gmail.com"}}, //to
		nil, //cc
		nil, //bcc
		"長いサブジェクト長いサブジェクト長いサブジェクト長いサブジェクト長いサブジェクト長いサブジェクト長いサブジェクト", //subject
		body.Bytes(), // body
	}
	mail.Send()
}

```
