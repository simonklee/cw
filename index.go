package main

type LinkIndex struct {
    monitor chan update
    index chan string
    store Store
}

func newLinkIndex(monitor chan update, store Store) *LinkIndex {
    l := &LinkIndex{
        monitor: monitor,
        index: make(chan string),
        store: store,
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
        
        p := initParse(string(data))
        p.next()
    }
}
