package main

import (
    "encoding/base64"
    "errors"
    "sync"
)

type Store struct {
    mu      sync.RWMutex
    entries map[string][]byte
    save    chan entry
}

type entry struct {
    id, url string
    data    []byte
}

func newStore() *Store {
    s := &Store{
        entries: make(map[string][]byte),
        save:    make(chan entry, 100),
    }

    go s.listen()
    return s
}

func (s *Store) listen() {
    for {
        select {
        case e := <-s.save:
            s.set(&e)

        }
    }
}

func (s *Store) get(id string) ([]byte, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if e, ok := s.entries[id]; ok {
        return e, nil
    }

    return nil, errors.New("not found")
}

func (s *Store) getByUrl(id, url string) ([]byte, error) {
    return s.get(s.key(url))
}

func (s *Store) set(e *entry) {
    s.mu.Lock()

    if e.id == "" {
        e.id = s.key(e.url)
    }

    defer s.mu.Unlock()
    s.entries[e.id] = e.data
}

func (s *Store) put(u string, data []byte) {
    e := entry{
        id:   s.key(u),
        url:  u,
        data: data,
    }

    s.save <- e
}

func (s *Store) key(u string) string {
    return base64.URLEncoding.EncodeToString([]byte(u))
}
