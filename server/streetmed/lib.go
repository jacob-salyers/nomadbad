package streetmed

import (
	"bufio"
	crypto "crypto/ed25519"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

type SearchMethod string
const LINEAR SearchMethod = "linear"
const TREE   SearchMethod = "tree"

type searchOpts struct { 
    SearchMethod SearchMethod
}

type set[T comparable] map[T]struct{}
func (this set[T]) Add(k T) { this[k] = struct{}{} }
func (this set[T]) Remove(k T) { delete(this, k) }
func (this set[T]) Contains(k T) bool { _, ok := this[k]; return ok }

func ToHTML(this set[string]) (string, error) {
    var sb strings.Builder
    sb.WriteString("<table><tr><th>File</th></tr>")
    
    for k := range this {
        sb.WriteString(fmt.Sprintf(`<tr><td><a href="/streetmed/files/%s.pdf" target="blank">%s</a></td></tr>`, k, html.EscapeString(string(k))))
    }
    sb.WriteString("</table>")

    return sb.String(), nil
}

func or[T comparable](in []set[T]) set[T] {
    ret := make(set[T])
    for _, el := range in {
        for k := range el {
            ret.Add(k)
        }
    }

    return ret
}

func and[T comparable](base *set[T], cmp []set[T]) error {
    var keysToRemove []T
    for k := range *base {
        for _, el := range cmp {
            if !el.Contains(k) {
                keysToRemove = append(keysToRemove, k)
                continue
            }
        }
    }

    for _, k := range keysToRemove { base.Remove(k) }

    return nil
}

func search(query query, opts searchOpts) (set[string], error) {

    var results []set[string]
    var err error
    switch opts.SearchMethod {
    case LINEAR:
        results, err = query.LinearSearch()
    case TREE:
        results, err = query.TreeSearch()
    default:
        err = errors.New("Invalid Method")
    }
    
    if err != nil { return nil, err }

    ret := or(results)
    if err := and(&ret, results); err != nil {
        return nil, err
    }

    return ret, nil
}

func scanCmdline(opts searchOpts) {
    scanner := bufio.NewScanner(os.Stdin)

    fmt.Println("Prompt?")
    for scanner.Scan(){
        input := scanner.Text()

        query, err := parseQuery(input)
        if err != nil { log.Fatal(err) }

        if input == "" { fmt.Println("Prompt?") ; continue }

        resultSet, err := search(query, opts)
        if err != nil { log.Fatal(err) }

        s := "files contain"
        if len(resultSet) == 1 {
            s = "file contains"
        }

        fmt.Printf("%d %s %s:\n", len(resultSet), s, query.CmdlineString())
        for k := range resultSet {
            fmt.Printf("\t%s.pdf\n", k)
        }
        fmt.Println("Prompt?")
    }
}

type tree [26]*node
type node struct {
    Leaf   set[int]
    Branch tree
}

func populateTree(labelIdx int, contents string, tree *tree) error {

    words := uniqueWords(contents)

    for _, word := range words {
        ptr := tree
        for _, c := range word {
            i := runeToIndex(c)
            // ignore non-lowercase ascii characters.
            if i < 0 { continue }
            if ptr[i]  == nil {
                tmp := node{ Leaf: make(set[int]) }
                ptr[i] = &tmp
            }
            ptr[i].Leaf.Add(labelIdx)
            ptr = &(ptr[i].Branch)
        }
    }

    return nil
}

func (this *tree) search(query string) set[string] {
    ret := make(set[string])
    var tmp *set[int]

    ptr := this
    for _, c := range query {
        i := runeToIndex(c)
        if i < 0 { continue }
        if ptr[i] == nil { return ret }

        tmp = &ptr[i].Leaf
        ptr = &ptr[i].Branch
    }
    
    for k := range *tmp { ret.Add(files[k]) }

    return ret
}

func runeToIndex(c rune) int {
    ret := int(c) - 97
    if ret < 0 || ret > 25 {
        return -1
    } else {
        return ret
    }
}

func uniqueWords(in string) []string {
    r, err := regexp.Compile("[^a-z]+")
    if err != nil { log.Fatal(err) }
    
    arr := r.Split(strings.ToLower(in), -1)
    sort.Strings(arr)
    var out []string
    for i := range arr {
        if i == 0 || arr[i] != arr[i-1] {
            out = append(out, arr[i])
        }
    }


    return out
}


// TODO (jacob): finish this
//               https://github.com/discord/discord-interactions-js/blob/main/src/index.ts
//               line 136:
func verifyDiscordSSLCert(w http.ResponseWriter, r *http.Request) bool {
    timestamp := r.Header.Get("X-Signature-Timestamp")
    signature := r.Header.Get("X-Signature-Ed25519")
    log.Printf("timestamp: '%s'", timestamp)
    log.Printf("signature: '%s'", signature)

    if timestamp == "" || signature == "" {
        return false
    }

    b, e := ioutil.ReadAll(r.Body)
    if e != nil {
        log.Println(e)
        return false
    }
    return crypto.Verify(_DISCORD_PUBLIC_KEY, b, []byte(signature))

}

func parseDiscordCred(b []byte) (string, string, string) {
    var id, pubkey, token string
    for _, line := range strings.Split(string(b), "\n") {
        arr := strings.Split(line, "\t")
        k := arr[0]
        v := arr[1]

        switch k {
        case "appid":
            id = v
        case "pubkey":
            pubkey = v
        case "token":
            token = v
        default:
            log.Fatalf("unrecognized discord cred key '%s'", k)
        }

    }

    return id, pubkey, token
}
