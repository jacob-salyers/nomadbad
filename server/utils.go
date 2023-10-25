package main

import (
	"fmt"
	"log"
	"net/http"
)


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
