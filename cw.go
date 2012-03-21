package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "sync"
    "strings"
    "encoding/base64"
)

var (
    logger = log.New(os.Stdout, "", 0)
    store *Store
)

type Store struct {
    mu      sync.RWMutex
    entries map[string] []byte
    save chan entry
}

type entry struct {
    id, url string
    data  []byte
}

type state struct {
    id string
    ok bool
}

func usage() {
    fmt.Fprintf(os.Stdout, "usage: cw url [url ...]\n")
    flag.PrintDefaults()
    os.Exit(0)
}

func main() {
    flag.Usage = usage
    flag.Parse()
    args := flag.Args()

    if len(args) == 0 {
        usage()
    } else {
        store = newStore()
        fetchUrls(args)

        <-store.save 
        println("main got e")
    }
}

func newStore() *Store {
    s := &Store{entries: make(map[string][]byte), save: make(chan entry, 100)}
    go s.listen()
    return s
}

func (s *Store) listen() {
    for {
        select {
        case e := <-s.save:
            println("save got e")
            s.set(&e)
        }
    }
}

func (s *Store) set(e *entry) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.entries[e.id] = e.data
}

func (s *Store) put(u string, data []byte) {
    e := entry{
        id: s.key(u),
        url: u, 
        data: data,
    }

    s.save<- e
}

func (s *Store) key(u string) string {
    return base64.URLEncoding.EncodeToString([]byte(u))
}

func fetch(u string) {
    i := strings.Index(u, "?")

    if i > 0 {
        u = u[:i]
    }

    res, err := http.Get(u)

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    debugResponse(res)

    data, err := ioutil.ReadAll(res.Body)
    defer res.Body.Close()

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    store.put(u, data)
}

func fetchUrls(urls []string) {
    for _, u := range urls {
        go fetch(u)
    }
}

func debugFetchState(fs *state) {
    logger.Printf("id: %s", fs.id)

    if fs.ok {
        logger.Printf("state: OK\n")
    } else {
        logger.Printf("state: ERR\n")
    }
}

func debugUrl(u *url.URL) {
    logger.Println("Host:", u.Host)
    logger.Println("Path:", u.Path)
    logger.Println("Request URI:", u.RequestURI())
    logger.Println("Scheme:", u.Scheme)
    logger.Println("Query:", u.RawQuery)
    logger.Println("Fragment:", u.Fragment)
}

func debugResponse(r *http.Response) {
    logger.Println("Status:", r.Status)
    logger.Println("StatusCode:", r.StatusCode)
    logger.Println("Proto:", r.Proto)
    logger.Println("Header:")

    for k, v := range r.Header {
        logger.Println("\t", k, ":", v)
    }
}

