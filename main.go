package main

import (
	"net/http"
	"log"
	"encoding/json"
	"os"
	"time"
)

func main() {
	// TODO (jacob): wrap the handler to do templating
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Print(http.ListenAndServe(":80", nil))
}

func logWrapper(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			log.Print(req.RemoteAddr)
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
