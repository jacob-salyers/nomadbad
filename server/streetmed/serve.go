package streetmed

import (
	"log"
	"net/http"
	"os"
	"strings"
)


var fileMap map[string]string
var files []string
var treeConst tree

func init() {
    // initialize file map
    fileMap = make(map[string]string)
    dirEntries, err := os.ReadDir("streetmed-data")
    if err != nil { log.Fatal(err) }

    for _, dirEntry := range dirEntries {
        arr := strings.Split(dirEntry.Name(),".")
        if arr[1] != "txt" { continue }

        tmp, err := os.ReadFile("streetmed-data/"+dirEntry.Name())
        if err != nil { log.Fatal(err) }

        fileMap[arr[0]] = strings.ToLower(string(tmp))
        files = append(files, arr[0])
    }

    // initialize tree
    for i := range files {
        err := populateTree(i, fileMap[files[i]], &treeConst)
        if err != nil { log.Fatal(err) }
    }

    http.Handle("GET /streetmed/api/protocol", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        queryString := r.URL.Query().Get("q")
        query, err := parseQuery(queryString)
        if err != nil { log.Println(err); w.WriteHeader(400); }

        set, err := search(query, searchOpts{TREE})
        if err != nil { log.Println(err); w.WriteHeader(500); }

        html, err := ToHTML(set)
        if err != nil { log.Println(err); w.WriteHeader(500); }

        w.Header().Add("Content-Type", "text/html")
        _, err = w.Write([]byte(html))
        if err != nil { log.Println(err); w.WriteHeader(500); }

        return
    }))

    http.Handle("/streetmed/", http.StripPrefix("/streetmed/", http.FileServer(http.Dir("streetmed-static"))))
    http.Handle("/streetmed/files/", http.StripPrefix("/streetmed/files/", http.FileServer(http.Dir("streetmed-files"))))
}
