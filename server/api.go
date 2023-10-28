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

	toPrint := fmt.Sprintf(`
From: %s %s
Reply-To: %s

%s
`,
	req.FormValue("first_name"), req.FormValue("last_name"),
	req.FormValue("email"), req.FormValue("message"));

	if _, err := f.Write([]byte(toPrint)); err != nil {
		log.Print(err)
		http.ServeFile(res, req, "static/submission_failure.html")
		return
	}

	if err := f.Close(); err != nil {
		log.Print(err)
		http.ServeFile(res, req, "static/submission_failure.html")
		return
	}

	http.ServeFile(res, req, "static/submission_successful.html")
}
