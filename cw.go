package main

import (
    "crypto/md5"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"
)

var (
    MaxCrawlers   = 1
    FetchInterval = time.Minute
)

type key string

type context struct {
    in      chan string
    store   Store
    monitor *Monitor
    index   *LinkIndex
    client  *http.Client
}

type request struct {
    id      key
    url     string
    store   Store
    monitor *Monitor
    index   chan key
    client  *http.Client
}

func newContext() *context {
    monitor := newMonitor()

    c := &context{
        in:      make(chan string, 1024),
        store:   NewMemoryStore(monitor.update),
        monitor: monitor,
        client:  &http.Client{},
    }

    c.index = NewLinkIndex(c.monitor.update, c.store, c.in)

    for i := 0; i < MaxCrawlers; i++ {
        go c.listen()
    }

    return c
}

func (c *context) listen() {
    for u := range c.in {
        c.Add(u, false)
    }
}

func (c *context) Add(u string, async bool) {
    u = c.normUrl(u)
    id := NewKey(u)

    if c.monitor.SetIfTime(id, StateIdle, StateFetch, FetchInterval) {
        r := c.newRequest(id, u)

        println()
        println("====== ", u, " ======")

        if async {
            go r.fetch()
        } else {
            r.fetch()
        }
    } else {
        //println("newly fetched", u)
        //c.in <- u
    }
}

func (c *context) normUrl(u string) string {
    n := strings.Index(u, "?")

    if n == -1 {
        n = len(u)
    }

    for u[n-1] == '/' && n >= 0 {
        n--
    }

    return u[:n]
}

func (c *context) newRequest(id key, url string) *request {
    r := new(request)
    r.url = url
    r.id = id
    r.store = c.store
    r.monitor = c.monitor
    r.index = c.index.in
    r.client = c.client
    return r
}

func (r *request) fetch() {
    var data []byte
    req, err := http.NewRequest("GET", r.url, nil)

    if err != nil {
        r.e(err)
        return
    }

    res, err := r.client.Do(req)

    if err != nil {
        r.e(err)
        return
    }

    //debugResponse(res)
    ct := res.Header.Get("Content-Type")

    if strings.Index(ct, "text/html") == -1 {
        return
    }

    data, err = ioutil.ReadAll(res.Body)

    if err != nil {
        r.e(err)
        return
    }

    defer res.Body.Close()

    if r.store.Set(&Entry{r.id, req.URL, data}) != nil {
        r.e(err)
        return
    }

    println("INDEXER", r.id)
    r.index <- r.id
    println("FINISHED", r.id)
    return
}

func (r *request) e(err error) {
    r.monitor.update <- update{r.id, StateError}
    fmt.Fprintln(os.Stderr, err)
}

func NewKey(url string) key {
    h := md5.New()
    h.Write([]byte(url))
    return key(fmt.Sprintf("%x", h.Sum(nil)))
}

func (k *key) String() string {
    return string(*k)
}

func (k *key) URL() *url.URL {
    return nil
}
