package main

import (
    "sync"
)

const (
    StateNone uint8 = 1 << iota
    StateIdle
    StateFetch
    StateStore
)

type State struct {
    id    string
    state uint8
}

type Monitor struct {
    mu    sync.RWMutex
    update chan State
    states map[string]uint8
}

func newMonitor() *Monitor {
    m := &Monitor{
        update: make(chan State, 100),
        states: make(map[string]uint8),
    }

    go m.listen()
    return m
}

func (m *Monitor) listen() {
    for {
        select {
        case state := <-m.update:
            m.Set(state.id, state.state)
        }
    }
}

func (m *Monitor) SetIf(id string, ifstate, state uint8) bool {
    m.mu.Lock()
    defer m.mu.Unlock()

    if oldstate, ok := m.states[id]; !ok || oldstate == ifstate {
        m.states[id] = state
        return true
    }

    return false
}

func (m *Monitor) Set(id string, state uint8) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.states[id] = state
    m.printState(id, state)
}

func (m *Monitor) Get(id string) uint8 {
    m.mu.RLock()
    defer m.mu.RUnlock()

    if state, ok := m.states[id]; ok {
        return state
    }

    return StateNone
}


func (m *Monitor) printState(id string, state uint8) {
    switch state {
    case StateIdle:
        println(id, "is_idle")
    case StateFetch:
        println(id, "is_fetch")
    case StateStore:
        println(id, "is_store")
    }
}
