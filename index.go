package main

import (
    "net/url"
)

type LinkIndex struct {
    monitor chan update
    in      chan key
    out     chan string
    store   Store
}

func NewLinkIndex(monitor chan update, store Store, out chan string) *LinkIndex {
    l := &LinkIndex{
        monitor: monitor,
        in:      make(chan key),
        store:   store,
        out:     out,
    }

    go l.listen()
    return l
}

func (l *LinkIndex) listen() {
    for id := range l.in {
        l.monitor <- update{id, StateIndex}
        e, err := l.store.Get(id)

        if err != nil {
            l.monitor <- update{id, StateError}
        }

        println("PARSING:", e.URL.String())
        urls := l.resolver(e)
        l.sorter(urls)
    }
}

func (l *LinkIndex) resolver(e *Entry) (urls []*url.URL) {
    raw, i := linkParse(e.Data, 0)
    //logger.Println(string(e.Data))

    for ; i > 0; raw, i = linkParse(e.Data, i) {
        u, err := url.Parse(string(raw))

        if err == nil {
            u = e.URL.ResolveReference(u)
            urls = append(urls, u)
        }
    }

    return
}

func (l *LinkIndex) sorter(urls []*url.URL) {
    for _, u := range urls {
        if u.Host == "simonklee.org" {
            //println("GOT", u.String())
            l.out <- u.String()
        } else {
            //println("IGNORE", u.String())
        }
    }
}
