package main

import (
    "errors"
    "os"
    "path/filepath"
    "sync"
    "net/url"
)

const (
    basepath = "/home/simon/src/github.com/simonz05/cw/store/"
)

type Store interface {
    Get(id key) (*Entry, error)
    Set(e *Entry) error
}

type Entry struct {
    Id   key
    URL  *url.URL
    Data []byte
}

type MemoryStore struct {
    monitor chan update

    // entries lock
    mu      sync.RWMutex
    entries map[key]*Entry
}

func NewMemoryStore(monitor chan update) *MemoryStore {
    s := &MemoryStore{
        entries: make(map[key]*Entry),
        monitor: monitor,
    }

    return s
}

func (s *MemoryStore) Set(e *Entry) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.monitor <- update{e.Id, StateStore}
    s.entries[e.Id] = e
    return nil
}

func (s *MemoryStore) Get(id key) (*Entry, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if e, ok := s.entries[id]; ok {
        return e, nil
    }

    return nil, errors.New("not found")
}

type FilesystemStore struct {
    monitor chan update
}

func (s *FilesystemStore) Set(e *Entry) error {
    p := s.path(e.Id)
    err := os.MkdirAll(filepath.Dir(p), 0754)

    if err != nil {
        logger.Println("mkdir:", err)
        return err
    }

    f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

    if err != nil {
        logger.Println("open:", err)
        return err
    }

    defer f.Close()
    n, err := f.Write(e.Data)
    // TODO: serialize entry to file

    if err != nil || n != len(e.Data) {
        logger.Println("write:", err)
        return err
    }

    return nil
}

func (s *FilesystemStore) path(id key) string {
    p := id[:2] + "/" + id[2:4] + "/" + id
    return string(basepath + p)
}

func (s *FilesystemStore) Get(id key) ([]byte, error) {
    return nil, nil
}
