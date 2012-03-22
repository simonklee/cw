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
    mu      sync.RWMutex
    entries map[string][]byte
    save    chan Entry
    update  chan State
}

type Entry struct {
    id   string
    data []byte
}

func newStore(update chan State) *Store {
    s := &Store{
        entries: make(map[string][]byte),
        save:    make(chan Entry, 100),
        update:  update,
    }

    go s.listen()
    return s
}

func (s *Store) listen() {
    for e := range s.save {
        s.set(&e)
        s.fsSet(e.id)
        s.update <- State{e.id, StateIdle}
    }
}

func (s *Store) set(e *Entry) {
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

func (s *Store) Put(id string, data []byte) {
    s.update <- State{id, StateStore}
    e := Entry{
        id:   id,
        data: data,
    }

    s.save <- e
}
