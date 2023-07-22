package main

import (
	"net/http"
	"log"
)

func main() {
	// TODO (jacob): wrap the handler to do templating
	http.Handle("/", logWrapper(http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/nomad-jiujitsu.com/fullchain.pem", "/etc/letsencrypt/live/nomad-jiujitsu.com/privkey.pem", nil))
}

func logWrapper(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			log.Print(req.RemoteAddr + " " + req.URL.Host + req.URL.Path + "?" + req.URL.RawQuery)
			wrappedHandler.ServeHTTP(res, req)
		})
}


/*
type blah struct {
	Id int
	Name string
}


func loadDatabase(path string, obj *blah) error {
	bytes, err := os.ReadFile(path)
	if (err != nil) {
		return err
	}

	err = json.Unmarshal(bytes, obj)
	if (err != nil) {
		return err
	}

	return nil
}
*/
