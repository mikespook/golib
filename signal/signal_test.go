package signal

import (
    "testing"
    "time"
    "os"
)

func TestSignalHandler(t *testing.T) {
    Bind(os.Interrupt, func() bool {return true})
    go func() {
        time.Sleep(time.Millisecond * 100)
        if err := Send(os.Getpid(), os.Interrupt); err != nil {
            t.Error(err)
        }
    }()
    s := Loop()
    t.Log(s)
}
