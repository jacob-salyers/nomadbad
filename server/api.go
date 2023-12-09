package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
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

	toPrint := fmt.Sprintf(`
From: %s %s
Reply-To: %s

%s
`,
	first, last, email, msg);

	if _, err := f.Write([]byte(toPrint)); err != nil {
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
