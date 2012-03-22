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
        store:   newStore(monitor.update),
        monitor: monitor,
    }

    return c
}

func (c *context) Add(u string) {
    id := c.key(u)

    if c.monitor.SetIf(id, StateIdle, StateFetch) {
        go c.fetch(id, u)
    }
}

func (c *context) fetch(id, u string) {
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
