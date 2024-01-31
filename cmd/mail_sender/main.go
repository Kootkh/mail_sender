package main

import (
	"flag"
	"fmt"
	"strings"

	"bytes"
	"encoding/base64"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Что бы мы могли задавать параметры в качестве флагов для нашего бинарника
type params struct {
	host                    string       `validate:"hostname"`
	name                    string       `validate:"omitempty,alphanumunicode"`
	from                    string       `validate:"omitempty,email"`
	fakeFrom                string       `validate:"omitempty,email"`
	recipients              []recipient  `validate:"omitempty,dive,required"`
	subject                 string       `validate:"required"`
	body                    string       `validate:"required"`
	attachments             []attachment `validate:"omitempty,dive,required"`
	charset                 string       `validate:"omitempty,oneof=UTF-8	Win-1251 CP-866	KOI-8R	ISO-8859-5"`
	contentType             string       `validate:"omitempty,oneof=EMAIL PHONE POST SMS"`
	encoding                string       `validate:"omitempty,oneof=EMAIL PHONE POST SMS"`
	contentTransferEncoding string       `validate:"omitempty,oneof=EMAIL PHONE POST SMS"`
	help                    *bool        `validate:"boolean"`
}

type recipient struct {
	email string `validate:"required,email"`
}

type attachment struct {
	attachment string `validate:"omitempty,file"`
}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

func init() {

	flag.String("host", "relay.obit.ru", "Host to send email")
	flag.String("name", "MAILER", "'helo' string")
	flag.String("from", "noreply@obit.ru", "Email address to send FROM")
	flag.String("fake-from", "", "[OPTIONAL] FAKE Email address to send FROM")
	flag.StringVar(&recipient, "to", "", "Email address to send TO")
	flag.String("subject", "", "Email SUBJECT")
	flag.String("body", "", "Email BODY")
	flag.StringVar(&attachment, "attach", "", "[OPTIONAL] Attachment (path to file)")
	flag.String("charset", "koi8-r", "[OPTIONAL] Email charset encoding")
	flag.String("contentType", "text/plain", "[OPTIONAL] Email 'content-type'")
	flag.String("encoding", "quoted-printable", "[OPTIONAL] Email encoding")
	flag.String("c-encoding", "8bit", "[OPTIONAL] Email 'content-transfer-encoding'")
	flag.Bool("help", false, "Print usage-help and exit")

}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

type Sender struct {
	auth smtp.Auth
}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

func (to *recipients) String() string {
	return fmt.Sprintf("%s", *to)
}

func (attach *attachments) String() string {
	return fmt.Sprintf("%s", *attach)
}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

func (to *recipients) Set(value string) error {
	if len(value) < 3 {
		return fmt.Errorf("to is too short")
	}

	*to = append(*to, value)

	return nil
}

func NewSender() *Sender {
	auth := smtp.PlainAuth("", username, password, host)
	return &Sender{auth}
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

func (attach *attachments) Set(value string) error {
	if len(value) < 3 {
		return fmt.Errorf("attach is too short")
	}

	*attach = append(*attach, value)

	return nil
}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

// main description of the Go function.
//
// No parameters.
// No return type.
func main() {

	sender := NewSender()
	m := NewMessage("Test", "Body message.")
	m.To = []string{"to@gmail.com"}
	m.CC = []string{"copy1@gmail.com", "copy2@gmail.com"}
	m.BCC = []string{"bc@gmail.com"}
	m.AttachFile("/path/to/file")
	fmt.Println(sender.Send(m))

	flag.Parse()

	fmt.Println("host:", *host)
	fmt.Println("name:", *name)
	fmt.Println("from:", *from)
	fmt.Println("fake-from:", *fakeFrom)
	fmt.Println("to:", to)
	fmt.Println("subject:", *subject)
	fmt.Println("body:", *body)
	fmt.Println("attachments:", attach)
	fmt.Println("charset:", *charset)
	fmt.Println("contentType:", *contentType)
	fmt.Println("encoding:", *encoding)
	fmt.Println("c-encoding:", *contentTransferEncoding)
	fmt.Println("help:", *help)

	/* --host="отправьте-пожалуйста.куда-нибудь.ауф"
	--name="pirozhok"
	--from="bober@kurwa.ауф"
	--fake-from="unicorn@kurwa.ауф"
	--to="victim@cheto.tam.chto-to"
	--to="ещё-одна-жертва@chet.tam.chto-to"
	--subject="Всё пропало. Клиент начал что-то подозревать. Высылайте аспирин в тюбиках!"
	--body="body-positive kurwa!"
	--attach="/паф/ту/филе"
	--attach="/овер/паф/ту/овер/филе"
	--charset="ЮТИФИ-ЭЙТ"
	--contentType="text/html"
	--encoding="узелковая-письменность"
	--c-encoding="буквосочетание"
	--help="враньё"
	*/

	/*
		--host="отправьте-пожалуйста.куда-нибудь.ауф" --name="pirozhok" --from="bober@kurwa.ауф" --fake-from="unicorn@kurwa.ауф" --to="victim@сhet.tam.chto-to" --to="ещё-одна-жертва@сhet.tam.chto-to" --subject="Всё пропало. Клиент начал что-то подозревать. Высылайте аспирин в тюбиках\!" --body="body-positive kurwa\!" --attach="/паф/ту/филе" --attach="/овер/паф/ту/овер/филе" --charset="ЮТИФИ-ЭЙТ" --contentType="text/html" --encoding="узелковая-письменность" --c-encoding="буквосочетание" --help="true"
	*/

	//sendMailSimple()
}

/* func sendMailSimple(subject string, body string, to {}string) {
	auth := smtp.PlainAuth(	// returns smtp.Auth object
		"",                   // identity (string)
	)

	headers := "MIME-Version: 1.0;\r\n" +
		"Content-Type: text/html;\r\n" +
		//"Content-Type: text/plain;
		"charset=\"" + charset +	"\";\r\n" +
		"Content-Transfer-Encoding: 8bit\r\n" +
		"From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n"

	//msg := "Subject: Cowabunga kurwa!\n Body-positive kurwa!" // message (string)
	msg := "Subject: " + subject + "\n" + headers + "\n\n" + html // message (string)

	err := smtp.SendMail(

		auth,											// smtp.Auth object

		to,												// to (string slice)
		[]byte(msg),              // message (slice of bytes)

	)

	if err != nil {
		fmt.Println(err)
	}
}
*/
