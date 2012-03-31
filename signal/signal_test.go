package signal

import (
    "testing"
    "time"
    "os"
)

func sendInterrupt(t *testing.T) {
    time.Sleep(time.Millisecond * 100)
    proc, err := os.FindProcess(os.Getpid())
    if err != nil {
        t.Error(err)
    }
    if err := proc.Signal(os.Interrupt); err != nil {
        t.Error(err)
    }
}

func TestSignalHandler(t *testing.T) {
    sh := NewHandler()
    sh.Bind(os.Interrupt, func() bool {return true})
    go sendInterrupt(t)
    s := sh.Loop()
    t.Log(s)
}
