package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"time"
	"net/http"
)

func apiLocalSignIn(res http.ResponseWriter, req *http.Request) {
	f, err := os.OpenFile("./data/signInQueue.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Print(err)
		http.ServeFile(res, req, "local/submission_failure.html")
		return
	}
	b, err := io.ReadAll(req.Body);
	if err != nil {
		log.Print(err)
		http.ServeFile(res, req, "local/submission_failure.html")
		return
	}

	t := time.Now().String()
	if _, err := f.Write([]byte(fmt.Sprintf("%s: %s\n", t, b))); err != nil {
		log.Print(err)
		http.ServeFile(res, req, "local/submission_failure.html")
		return
	}

	if err := f.Close(); err != nil {
		log.Print(err)
		http.ServeFile(res, req, "local/submission_failure.html")
		return
	}

	http.ServeFile(res, req, "local/submission_successful.html")
}

func apiLocalSignUp(res http.ResponseWriter, req *http.Request) {
	f, err := os.OpenFile("./data/signUpQueue.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Print(err)
		http.ServeFile(res, req, "local/submission_failure.html")
		return
	}

	b, err := io.ReadAll(req.Body);
	if  err != nil {
		log.Print(err)
		http.ServeFile(res, req, "local/submission_failure.html")
		return
	}

	t := time.Now().String()
	if _, err := f.Write([]byte(fmt.Sprintf("%s: %s\n", t, b))); err != nil {
		log.Print(err)
		http.ServeFile(res, req, "local/submission_failure.html")
		return
	}

	if err := f.Close(); err != nil {
		log.Print(err)
		http.ServeFile(res, req, "local/submission_failure.html")
		return
	}

	http.ServeFile(res, req, "local/submission_successful.html")
}
