package signal

import (
    "testing"
    "time"
    "os"
)

func TestSignalHandler(t *testing.T) {
    sh := NewHandler()
    sh.Bind(os.Interrupt, func() bool {return true})
    go func() {
        time.Sleep(time.Millisecond * 100)
        if err := sh.Send(os.Getpid(), os.Interrupt); err != nil {
            t.Error(err)
        }
    }()
    s := sh.Loop()
    t.Log(s)
}
