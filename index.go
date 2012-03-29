package main

import (
    "net/url"
    "regexp"
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

        //println("PARSING:", e.URL.String())
        urls := l.resolver(e)
        l.localSorter(urls)
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

var (
    // https://github.com/{user}/{project}/blob/{branch}/{filename}
    //
    // https://github.com/(simonz05/godis)/blob/(stable/conn.go)
    //    =
    // https://github.com/(simonz05/godis)/raw/(stable/conn.go)
    //    =
    // https://raw.github.com/(simonz05/godis)/(stable/conn.go)
    githubResRegex = regexp.MustCompile("^/(.*/.*)/blob/(.*)$")
    githubResRepl = "/$1/$2"

    //https://github.com/{user}/{project}/blob/{branch}/{filename}
    localResRegex = regexp.MustCompile("^/blob/(.*)")
    localResRepl = "/raw/$1"
)

func (l *LinkIndex) localSorter(urls []*url.URL) {
    for _, u := range urls {
        switch {
        case u.Host != "localhost":
            // ignore
        case localResRegex.MatchString(u.Path):
            u.Path = localResRegex.ReplaceAllString(u.Path, localResRepl)
            // TODO save in redis resource queue
        default:
            // TODO save to redis link queue
            // <- is a deadlock because links found will always be greater than links consumed(crawled)
            l.out <- u.String()
        } 
    }
}
