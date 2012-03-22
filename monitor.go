package main

import ()

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
            m.setState(&state)
        }
    }
}

func (m *Monitor) setState(state *State) {
    m.states[state.id] = state.state
    m.printState(state.id)
}

func (m *Monitor) printState(id string) {
    s := m.states[id]

    switch s {
    case StateIdle:
        println(id, "is_idle")
    case StateFetch:
        println(id, "is_fetch")
    case StateStore:
        println(id, "is_store")
    }
}
