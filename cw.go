package main

import (
    "crypto/md5"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
)

type context struct {
    in      chan string
    store   *Store
    monitor *Monitor
}

func newContext() *context {
    monitor := newMonitor()
    c := &context{
        in:      make(chan string),
        store:   newStore(monitor.update),
        monitor: monitor,
    }

    go c.listen()
    return c
}

func (c *context) Add(u string) {
    c.in <- u
}

func (c *context) listen() {
    for {
        select {
        case u := <-c.in:
            id := c.key(u)
            c.monitor.update <- State{id, StateIdle}
            go c.fetch(u, id)
        }
    }
}

func (c *context) fetch(u, id string) {
    c.monitor.update <- State{id, StateFetch}
    i := strings.Index(u, "?")

    if i > 0 {
        u = u[:i]
    }

    res, err := http.Get(u)

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    //debugResponse(res)

    data, err := ioutil.ReadAll(res.Body)
    defer res.Body.Close()

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    c.store.Put(id, data)
}

func (c *context) key(u string) string {
    h := md5.New()
    h.Write([]byte(u))
    return fmt.Sprintf("%x", h.Sum(nil))
}
