package main

import (
    "errors"
    "os"
    "path/filepath"
    "sync"
)

const (
    basepath = "/home/simon/src/github.com/simonz05/cw/store/"
)

type Store struct {
    Save    chan entry
    monitor chan update

    // entries lock
    mu      sync.RWMutex
    entries map[string][]byte
}

type entry struct {
    id   string
    data []byte
}

func newStore(monitor chan update) *Store {
    s := &Store{
        entries: make(map[string][]byte),
        monitor: monitor,
        Save:    make(chan entry, 100),
    }

    go s.listen()
    return s
}

func (s *Store) listen() {
    for e := range s.Save {
        s.monitor <- update{e.id, StateStore}
        s.set(&e)
        s.fsSet(e.id)
        s.monitor <- update{e.id, StateIdle}
    }
}

func (s *Store) set(e *entry) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.entries[e.id] = e.data
}

func (s *Store) fsSet(id string) error {
    data, err := s.Get(id)

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

    defer f.Close()

    if n, err := f.Write(data); err != nil || n != len(data) {
        return err
    }

    return nil
}

func (s *Store) path(id string) string {
    p := id[:2] + "/" + id[2:4] + "/" + id
    return basepath + p
}

func (s *Store) Get(id string) ([]byte, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if e, ok := s.entries[id]; ok {
        return e, nil
    }

    return nil, errors.New("not found")
}
