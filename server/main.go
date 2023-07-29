package main

import (
	"net/http"
	"log"
	"flag"
)

func main() {
	useSSL := flag.Bool("s", false, "Toggles http/https")
	flag.Parse()
	// TODO (jacob): wrap the handler to do templating

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
			wrappedHandler.ServeHTTP(res, req)
		})
}

func redirectToHTTPS() {
	log.Print("Starting http redirector")
	httpServer := http.Server{
		Addr: ":80",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Print("in redirector")
			u := r.URL
			u.Scheme = "https"
			log.Print(u.String())
			http.Redirect(w,r,u.String(), http.StatusMovedPermanently)
		}),
	}
	log.Fatal(httpServer.ListenAndServe())
}
