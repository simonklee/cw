package main

import (
    "testing"
)

type monitorTest struct {
    url string
    ok  bool
}

var monitorUrls = []monitorTest{
    {"http://foo.com", true},
    {"http://foo.com/", true},
    {"http://foo.com/", false},
    {"http://foo.com/?foo=bar", true},
}

func TestSetIf(t *testing.T) {
    m := newMonitor()

    for _, test := range monitorUrls {
        if ok := m.SetIf(NewKey(test.url), StateFetch, StateIdle); ok != test.ok {
            error_(t, test.ok, ok, nil)
        } else {
            t.Log(test.url, test.ok, ok)
        }
    }

    close(m.update)
}
