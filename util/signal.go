package util

import (
    "os"
    "os/signal"
)

// return true if need to break
type SignalCallback func() bool

type SignalHandler struct {
    schan chan os.Signal
    cb map[os.Signal]SignalCallback
}

func NewSignalHandler() (sh *SignalHandler) {
    sh = &SignalHandler{make(chan os.Signal, 1), make(map[os.Signal]SignalCallback, 5)}
    signal.Notify(sh.schan, os.Interrupt, os.Kill)
    return
}

func (sh *SignalHandler)Bind(s os.Signal, cb SignalCallback) {
    sh.cb[s] = cb
}

func (sh *SignalHandler)Loop() os.Signal {
    for {
        s := <-sh.schan
        f, ok := sh.cb[s]
        if ok {
            f()
            return s
        }
    }
    return nil
}
