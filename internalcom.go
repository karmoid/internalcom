package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/go-mail/mail"
)

const SMTP_USER_ENV string = "SMTP_USER"
const SMTP_PWD_ENV string = "SMTP_PWD"
const LOGO_FILENAME string = "logo.jpg"

const SMTP_ENDPOINT string = "smtp.office365.com"

type Mail struct {
	ToName   string
	ToAddr   string
	FromName string
	FromAddr string
	Subject  string
	Body     string
	Logofile string
	Port     int
}

func mailer(ml Mail, logo bool) error {

	smtp_user := strings.TrimSpace(os.Getenv(SMTP_USER_ENV))
	if smtp_user == "" {
		log.Fatal("Deploy: Unable to retrieve SMTP user. Make sure environment sets: " + SMTP_USER_ENV)
	}

	smtp_pwd := strings.TrimSpace(os.Getenv(SMTP_PWD_ENV))
	if smtp_pwd == "" {
		log.Fatal("Deploy: Unable to retrieve SMTP password. Make sure environment sets: " + SMTP_PWD_ENV)
	}

	if ml.FromAddr == "" {
		log.Fatal("Run: Unable to retrieve FromAddr. Make sure arg sets: From")
	}

	if ml.ToAddr == "" {
		log.Fatal("Run: Unable to retrieve ToAddr. Make sure arg sets: To")
	}

	if ml.Subject == "" {
		log.Fatal("Run: Unable to retrieve Subject. Make sure arg sets: Subject")
	}

	if ml.Body == "" {
		log.Fatal("Run: Unable to retrieve Body. Make sure arg sets: Body")
	}

	mailer := SMTP_ENDPOINT
	m := mail.NewMessage()
	if ml.FromName == "" {
		m.SetHeader("From", ml.FromAddr)
	} else {
		m.SetHeader("From", m.FormatAddress(ml.FromAddr, ml.FromName))
	}
	m.SetHeader("From", m.FormatAddress(ml.FromAddr, ml.FromName))

	recipients := strings.SplitN(ml.ToAddr, ";", -1)
	addresses := make([]string, len(recipients))
	for i, recipient := range recipients {
		addresses[i] = m.FormatAddress(recipient, "")
		// log.Println("%d - %s", i, recipient)
	}
	// log.Println("user:[" + smtp_user + "] pwd:[" + smtp_pwd + "]")

	m.SetHeader("To", addresses...)
	m.SetHeader("Subject", ml.Subject)
	if logo {
		m.Embed(ml.Logofile)
		m.SetBody("text/html", ml.Body+"<br/><img src=\"cid:"+ml.Logofile+"\" alt=\"Logo\" align=\"right\"/>")
	} else {
		m.SetBody("text/html", ml.Body)
	}

	d := mail.NewDialer(mailer, ml.Port, smtp_user, smtp_pwd)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// Send the email
	for index := 0; index < 10; index++ {
		// log.Println("trying a send ", index)
		err := d.DialAndSend(m)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "timeout") {
			log.Println("timeout - retry ", index+1)
			if index < 8 {
				continue
			}
		}
		return err
	}
	return nil
}

// main - Entry pont for Internal Communication sender program
// Goal is to use O365 services to send email with Embbeded logo
//
// V1.0 - Initial version
// V1.1 - Port and Logofile args added, "No need to embed Logo if nolog choose" request
func main() {
	log.Println("internalcom - Internal Communication - Email through O365 - C.m. V1.1")

	toAddrPtr := flag.String("to", "", "To email address")
	fromNamePtr := flag.String("sender", "", "Sender Name (First & Last)")
	fromAddrPtr := flag.String("from", "", "From email address")
	subjectPtr := flag.String("subject", "", "Subject of email")
	bodyPtr := flag.String("body", "", "Body of email")
	logofilePtr := flag.String("logofile", LOGO_FILENAME, "Logo filename (jpg | png)")
	logoPtr := flag.Bool("logo", true, "Put logo at end of email")
	portPtr := flag.Int("port", 587, "smtp port")

	flag.Parse()

	m := Mail{ToAddr: *toAddrPtr,
		FromName: *fromNamePtr,
		FromAddr: *fromAddrPtr,
		Subject:  *subjectPtr,
		Body:     *bodyPtr,
		Port:     *portPtr,
		Logofile: *logofilePtr}

	err := mailer(m, *logoPtr)
	if err != nil {
		log.Fatal(err)
	}
}