package util

import (
    "testing"
)

const (
    PIDFILE = "./test.pid"
)

func TestPidFile(t *testing.T) {
    pf, err := NewPidFile(PIDFILE)
    if err != nil {
        t.Error(err)
    }

    if err := pf.Close(); err != nil {
        t.Error(err)
    }
}
