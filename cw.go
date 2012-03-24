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

type request struct {
    id      string
    url     string
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
    u = c.normUrl(u)
    id := c.id(u)

    if c.monitor.SetIf(id, StateIdle, StateFetch) {
        r := c.newRequest(id, u)
        go r.fetch()
    }
}

func (c *context) id(u string) string {
    h := md5.New()
    h.Write([]byte(u))
    return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *context) normUrl(u string) string {
    i := strings.Index(u, "?")

    if i == -1 {
        i = len(u)
    }

    return u[:i]
}

func (c *context) newRequest(url, id string) *request {
    r := new(request)
    r.url = url
    r.id = id
    r.store = c.store
    r.monitor = c.monitor
    return r
}

func (r *request) fetch() {
    var data []byte
    res, err := http.Get(r.url)

    if err != nil {
        goto Error
    }

    //debugResponse(res)
    data, err = ioutil.ReadAll(res.Body)

    if err != nil {
        goto Error
    }

    defer res.Body.Close()

    r.store.Save <- entry{r.id, data}
    return
Error:
    r.monitor.update <- update{r.id, StateError}
    fmt.Fprintln(os.Stderr, err)
}
