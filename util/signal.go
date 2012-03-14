package util

import (
    "os"
    "os/signal"
)

// return true if need to break
type SignalCallback func() bool

type SignalHandler struct {
    schan chan os.Signal
    icb SignalCallback
    kcb SignalCallback
}

func NewSignalHandler(icb, kcb SignalCallback) (sh *SignalHandler) {
    sh = &SignalHandler{make(chan os.Signal, 1), icb, kcb}
    signal.Notify(sh.schan, os.Interrupt, os.Kill)
    return
}

func (sh *SignalHandler)Loop() os.Signal {
    for {
        s := <-sh.schan
        switch(s) {
            case os.Interrupt:
                if (sh.icb()) {
                    return s
                }
            case os.Kill:
                if (sh.kcb()) {
                    return s
                }
        }
    }
    return nil
}
