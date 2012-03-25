package main

import (
    "net/url"
)

type LinkIndex struct {
    monitor chan update
    index   chan key
    store   Store
}

func NewLinkIndex(monitor chan update, store Store) *LinkIndex {
    l := &LinkIndex{
        monitor: monitor,
        index:   make(chan key),
        store:   store,
    }

    go l.listen()
    return l
}

func (l *LinkIndex) listen() {
    for id := range l.index {
        l.monitor <- update{id, StateIndex}
        data, err := l.store.Get(id)

        if err != nil {
            l.monitor <- update{id, StateError}
        }

        urlResolver(data)
    }
}

func urlResolver(e *Entry) (urls []*url.URL) {
    raw, i := linkParse(e.Data, 0); 
    //logger.Println(string(e.Data))

    for ; i > 0; raw, i = linkParse(e.Data, i) {
        u, err := url.Parse(string(raw)); 

        if err == nil {
            u = e.URL.ResolveReference(u)
            urls = append(urls, u)
        }
    }

    return
}

func absUrl(u *url.URL, domain string) (*url.URL) {
    debugUrl(u)
    return u 
}
