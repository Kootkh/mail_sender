package main

import (
	"flag"
	"fmt"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Что бы мы могли задавать параметры в качестве флагов для нашего бинарника
type params struct {
	host                    string `validate:"hostname"`
	name                    string
	from                    string        `validate:"email"`
	fakeFrom                string        `validate:"email"`
	recipients              []*recipient  `validate:"required,dive,required"`
	subject                 string        `validate:"required"`
	body                    string        `validate:"required"`
	attachments             []*attachment `validate:"omitempty,dive,required"`
	charset                 string
	contentType             string
	encoding                string
	contentTransferEncoding string
	help                    bool `validate:"boolean"`
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
	flag.StringVar(&to, "to", "Email address to send TO")
	flag.String("subject", "", "Email SUBJECT")
	flag.String("body", "", "Email BODY")
	flag.StringVar(&attach, "attach", "[OPTIONAL] Attachment (path to file)")
	flag.String("charset", "koi8-r", "[OPTIONAL] Email charset encoding")
	flag.String("contentType", "text/plain", "[OPTIONAL] Email 'content-type'")
	flag.String("encoding", "quoted-printable", "[OPTIONAL] Email encoding")
	flag.String("c-encoding", "8bit", "[OPTIONAL] Email 'content-transfer-encoding'")
	flag.Bool("help", false, "Print usage-help and exit")

}

func (to *recipients) String() string {
	return fmt.Sprintf("%s", *to)
}

func (attach *attachments) String() string {
	return fmt.Sprintf("%s", *attach)
}

func (to *recipients) Set(value string) error {
	if len(value) < 3 {
		return fmt.Errorf("to is too short")
	}

	*to = append(*to, value)

	return nil
}

func (attach *attachments) Set(value string) error {
	if len(value) < 3 {
		return fmt.Errorf("attach is too short")
	}

	*attach = append(*attach, value)

	return nil
}

// main description of the Go function.
//
// No parameters.
// No return type.
func main() {

	var to recipients
	var attach attachments

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
