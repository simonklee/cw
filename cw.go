package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
)

type monitor struct {
    in, out chan string
}

func newMonitor() *monitor {
    m := &monitor{
        in:  make(chan string, 100),
        out: make(chan string),
    }

    go m.listen()

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
