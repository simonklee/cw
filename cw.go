package main

import (
    "crypto/md5"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
)

type key string

type context struct {
    in      chan string
    store   Store
    monitor *Monitor
    index   *LinkIndex
}

type request struct {
    id      key
    url     string
    store   Store
    monitor *Monitor
    index   chan key
}

func newContext() *context {
    monitor := newMonitor()
    c := &context{
        store:   NewMemoryStore(monitor.update),
        monitor: monitor,
    }
    c.index = NewLinkIndex(c.monitor.update, c.store)
    return c
}

func (c *context) Add(u string) {
    u = c.normUrl(u)
    id := NewKey(u)

    if c.monitor.SetIf(id, StateIdle, StateFetch) {
        r := c.newRequest(id, u)
        go r.fetch()
    }
}

func (c *context) normUrl(u string) string {
    i := strings.Index(u, "?")

    if i == -1 {
        i = len(u)
    }

    return u[:i]
}

func (c *context) newRequest(id key, url string) *request {
    r := new(request)
    r.url = url
    r.id = id
    r.store = c.store
    r.monitor = c.monitor
    r.index = c.index.index
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

    if r.store.Set(&Entry{r.id, data}) != nil {
        goto Error
    }

    r.index <- r.id
    return
Error:
    r.monitor.update <- update{r.id, StateError}
    fmt.Fprintln(os.Stderr, err)
}

func NewKey(url string) key {
    h := md5.New()
    h.Write([]byte(url))
    return key(fmt.Sprintf("%x", h.Sum(nil)))
}

func (k key) String() string {
    return string(k)
}
