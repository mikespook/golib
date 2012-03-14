package util

import (
    "testing"
)

func TestSignalHandler(t *testing.T) {
    sh := NewSignalHandler(func() bool {return true}, func() bool {return false})
    s := sh.Loop()
    t.Log(s)
}
