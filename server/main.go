package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	useSSL := flag.Bool("s", false, "Toggles http/https")
	local := flag.Bool("l", false, "Toggles local-only pages")
	adhoc := flag.Bool("a", false, "Toggles adhoc code")
	flag.Parse()

    if *adhoc {

        tmp := MailGunInput {
            ReplyTo: "Jacob Salyers <kswnin@gmail.com>",
            From: "Nomad Form <mailgun@mg.nomad-jiujitsu.com>",
            To: "caravancollective@outlook.com",
            Subject: "Go Test 1",
            Body: "This is some text!",
        }
        mailgun(tmp)
        os.Exit(0)
    }

	if *local {
		http.Handle("/api-local/sign-in", http.HandlerFunc(apiLocalSignIn))
		http.Handle("/api-local/sign-up", http.HandlerFunc(apiLocalSignUp))
		http.Handle("/local/", http.StripPrefix("/local/", http.FileServer(http.Dir("local"))))
	}

	http.Handle("/api/submit", logWrapper(http.HandlerFunc(apiSubmit)))
	http.Handle("/", logWrapper(http.FileServer(http.Dir("static"))))

	if *useSSL {
		log.Print("Starting on port 443, with redirect from 80")
		go redirectToHTTPS()
		log.Fatal(http.ListenAndServeTLS(
			":443",
			"/etc/letsencrypt/live/nomad-jiujitsu.com/fullchain.pem",
			"/etc/letsencrypt/live/nomad-jiujitsu.com/privkey.pem",
			nil))
	} else {
		log.Print("Starting on port 8000")
		log.Fatal(http.ListenAndServe(":8000", nil))
	}
}
