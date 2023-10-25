package main

import (
	"fmt"
	"os"
	"log"
	"flag"
	"net/http"
)

func main() {
	useSSL := flag.Bool("s", false, "Toggles http/https")
	flag.Parse()
	// TODO (jacob): wrap the handler to do templating

	http.Handle("/api/submit", logWrapper(http.HandlerFunc(apiSubmit)))

	http.Handle("/", logWrapper(http.FileServer(http.Dir("static"))))

	if *useSSL {
		log.Print("Starting on port 443, with redirect from 80")
		go redirectToHTTPS()
		log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/nomad-jiujitsu.com/fullchain.pem", "/etc/letsencrypt/live/nomad-jiujitsu.com/privkey.pem", nil))
	} else {
		log.Print("Starting on port 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}

func logWrapper(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
            log.Print(req.RemoteAddr + " " + req.URL.Host + req.URL.Path + "?" + req.URL.RawQuery)
            fmt.Println(req.URL.Path)

			wrappedHandler.ServeHTTP(res, req)
		})
}

func redirectHelper(res http.ResponseWriter, req *http.Request) {
	log.Print("Redirecting...")
	http.Redirect(res, req, "https://nomad-jiujitsu.com" + req.RequestURI, http.StatusMovedPermanently)
}

func redirectToHTTPS() {
	log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(redirectHelper)))
}

func apiSubmit(res http.ResponseWriter, req *http.Request) {
	var message = req.FormValue("message")
	log.Println(message)
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
