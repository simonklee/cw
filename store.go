package main

import (
    "crypto/md5"
    "errors"
    "sync"
    "fmt"
    "path/filepath"
    "os"
    "bufio"
)

const (
    basepath = "/home/simon/src/github.com/simonz05/cw/store/"
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
    for e := range s.save {
        s.set(&e)
        s.fsSet(e.id)
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

func (s *Store) getByUrl(u string) ([]byte, error) {
    return s.get(s.key(u))
}

func (s *Store) put(u string, data []byte) {
    e := entry{
        id:   s.key(u),
        url:  u,
        data: data,
    }

    s.save <- e
}

func (s *Store) set(e *entry) {
    s.mu.Lock()

    if e.id == "" {
        e.id = s.key(e.url)
    }

    defer s.mu.Unlock()
    s.entries[e.id] = e.data
}

func (s *Store) fsSet(id string) error {
    data, err := s.get(id)

    if err != nil {
        return err
    }

    p := s.path(id)

    if err := os.MkdirAll(filepath.Dir(p), 0754); err != nil {
        logger.Println("mkdir:", err)
        return err
    }

    f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

    if err != nil {
        logger.Println("open:", err)
        return err
    }

    b := bufio.NewWriter(f)
    defer f.Close()
    defer b.Flush()

    if n, err := b.Write(data); err != nil || n != len(data) {
        return err
    }

    return nil
}

func (s *Store) path(id string) string {
    p := id[:2] + "/" + id[2:4] + "/" + id
    return basepath + p
}

func (s *Store) key(u string) string {
    h := md5.New()
    h.Write([]byte(u))
    return fmt.Sprintf("%x", h.Sum(nil))
}
