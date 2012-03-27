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
    last   time.Time
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

    s.last = time.Now()
    s.status = status
    m.states[s.id] = s
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

func (m *Monitor) SetIfTime(id key, ifstatus, status flags, interval time.Duration) bool {
    m.mu.Lock()
    defer m.mu.Unlock()
    s, ok := m.states[id]

    if !ok || (s.status == ifstatus && time.Since(s.last) > interval) {
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
