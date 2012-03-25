package main

import (
    "sync"
    "time"
)

type flags uint32

const (
    StateNone flags = iota
    StateIdle
    StateFetch
    StateStore
    StateIndex
    StateError
)

type state struct {
    id     key
    status flags
    last   int64
}

type update struct {
    id     key
    status flags
}

type Monitor struct {
    update chan update

    mu     sync.RWMutex
    states map[key]state
}

func newMonitor() *Monitor {
    m := &Monitor{
        update: make(chan update, 100),
        states: make(map[key]state),
    }

    go m.listen()
    return m
}

func (m *Monitor) listen() {
    for u := range m.update {
        m.mu.Lock()
        m.set(u.id, u.status)
        m.mu.Unlock()
    }
}

func (m *Monitor) set(id key, status flags) {
    s, ok := m.states[id]

    if !ok {
        s = state{}
        m.states[id] = s
    }

    s.last = time.Now().Unix()
    s.status = status
    m.states[s.id] = s
    m.printState(s.id, s.status)
}

func (m *Monitor) SetIf(id key, ifstatus, status flags) bool {
    m.mu.Lock()
    defer m.mu.Unlock()

    if s, ok := m.states[id]; !ok || s.status == ifstatus {
        m.set(id, status)
        return true
    }

    return false
}

func (m *Monitor) SetIfTime(id key, ifstatus, status flags, d int64) bool {
    m.mu.Lock()
    defer m.mu.Unlock()
    s, ok := m.states[id]

    if !ok || (s.status == ifstatus && time.Now().Unix()-s.last > d) {
        m.set(id, status)
        return true
    }

    return false
}

func (m *Monitor) Get(id key) flags {
    m.mu.RLock()
    defer m.mu.RUnlock()

    if s, ok := m.states[id]; ok {
        return s.status
    }

    return StateNone
}

func (m *Monitor) printState(id key, status flags) {
    switch status {
    case StateIdle:
        println(id, "is_idle")
    case StateFetch:
        println(id, "is_fetch")
    case StateStore:
        println(id, "is_store")
    case StateError:
        println(id, "is_error")
    case StateIndex:
        println(id, "is_index")
    case StateNone:
        println(id, "is_none")
    }
}
