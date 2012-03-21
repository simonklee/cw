package main

import (
    "encoding/base64"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "strings"
    "errors"
    "sync"
)

var (
    logger = log.New(os.Stdout, "", 0)
)

type Store struct {
    mu      sync.RWMutex
    entries map[string][]byte
    save    chan entry
}

type entry struct {
    id, url string
    data    []byte
}

type backend struct {
    save    chan<- entry
    in, out chan string
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
        b := newBackend()

        for _, u := range args {
            b.in <- u
        }

        b.listen()
    }
}

func newStore() *Store {
    s := &Store{
        entries: make(map[string][]byte),
        save:    make(chan entry, 100),
    }

    go s.listen()
    return s
}

func (s *Store) listen() {
    for {
        select {
        case e := <-s.save:
            s.set(&e)
        }
    }
}

func (s *Store) get(id string) ([]byte, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if e, ok := s.entries[id]; ok {
        return e, nil
    }

    return nil, errors.New("url not found")
}

func (s *Store) getByUrl(id, url string) ([]byte, error) {
    return s.get(s.key(url))
}

func (s *Store) set(e *entry) {
    s.mu.Lock()

    if e.id == "" {
        e.id = s.key(e.url)
    }

    defer s.mu.Unlock()
    s.entries[e.id] = e.data
}

func (s *Store) put(u string, data []byte) {
    e := entry{
        id:   s.key(u),
        url:  u,
        data: data,
    }

    s.save <- e
}

func (s *Store) key(u string) string {
    return base64.URLEncoding.EncodeToString([]byte(u))
}

func newBackend() *backend {
    s := newStore()

    b := &backend{
        save: s.save,
        in:   make(chan string, 100),
        out:  make(chan string),
    }

    return b
}

func (b *backend) listen() {
    for {
        select {
        case u := <-b.in:
            println("--> in ", u)
            go fetch(u, b.save)
        case u := <-b.out:
            println("<-- out ", u)
        }
    }
}

func fetch(u string, save chan<- entry) {
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

    save <- entry{url: u, data: data}
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
