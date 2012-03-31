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
