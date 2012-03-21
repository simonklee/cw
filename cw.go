package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
)

var (
    logger = log.New(os.Stdout, "", 0)
    store *Store
)

type monitor struct {
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
        m := newMonitor()
        store = newStore()

        for _, u := range args {
            m.in <- u
        }

        m.listen()
    }
}

func newMonitor() *monitor {
    m := &monitor{
        in:   make(chan string, 100),
        out:  make(chan string),
    }

    return m
}

func (m *monitor) listen() {
    for {
        select {
        case u := <-m.in:
            println("--> in ", u)
            go fetch(u)
        case u := <-m.out:
            println("<-- out ", u)
        }
    }
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

    store.save <- entry{url: u, data: data}
}
