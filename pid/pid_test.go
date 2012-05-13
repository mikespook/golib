package pid

import (
    "testing"
)

const (
    PIDFILE = "./test.pid"
)

func TestPidFile(t *testing.T) {
    pf, err := New(PIDFILE)
    if err != nil {
        t.Error(err)
    }

    if err := pf.Close(); err != nil {
        t.Error(err)
    }
}

func TestProcessExists(t *testing.T) {
    pf1, err := New(PIDFILE)
    if err != nil {
        t.Error(err)
    }
    defer pf1.Close()
    pf2, err := New(PIDFILE)
    if err != ProcessExists {
        t.Error("The process should exist.")
    }
    if pf2 != nil {
        defer pf2.Close()
    }
}
