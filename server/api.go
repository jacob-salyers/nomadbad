package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func apiSubmit(res http.ResponseWriter, req *http.Request) {

	f, err := os.OpenFile("./form_submissions.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err)
		http.ServeFile(res, req, "static/submission_failure.html")
		return
	}

    first := req.FormValue("first_name")
    last := req.FormValue("last_name")
    email := req.FormValue("email")
    msg := req.FormValue("message")

    now := time.Now()
    var sb strings.Builder
    sb.WriteString(now.String())
    sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("From: %s %s\n", first, last))
    sb.WriteString(fmt.Sprintf("Reply-To: %s\n", email))
    sb.WriteString("\t> ")
    i := 0
    for _, el := range strings.Split(msg, "") {
        if (i >= 70 && el == " ") || el == "\n" {
            i = 0
            sb.WriteString("\n\t> ")
        } else {
            sb.WriteString(el)
            i += 1
        }
    }
    sb.WriteString("\n\n")

	if _, err := f.Write([]byte(sb.String())); err != nil {
		log.Print(err)
		http.ServeFile(res, req, "static/submission_failure.html")
		return
	}

    err = mailgun(MailGunInput {
        ReplyTo: fmt.Sprintf( "%s %s <%s>", 
            first,
            last,
            email),
        From: "Nomad Form <mailgun@mg.nomad-jiujitsu.com>",
        To: "caravancollective@outlook.com",
        Subject: fmt.Sprintf("Form Submission from %s %s", 
            first,
            last),
        Body: msg,
        })

    if err != nil {
        log.Println(err)
    }

	if err := f.Close(); err != nil {
		log.Print(err)
		http.ServeFile(res, req, "static/submission_failure.html")
		return
	}

	http.ServeFile(res, req, "static/submission_successful.html")
}
