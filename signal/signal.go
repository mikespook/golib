package signal

import (
    "os"
    S "os/signal"
)

// return true if need to break
type Callback func() bool

type Handler struct {
    schan chan os.Signal
    cb map[os.Signal]Callback
}

func NewHandler() (sh *Handler) {
    sh = &Handler{make(chan os.Signal, 1), make(map[os.Signal]Callback, 5)}
    S.Notify(sh.schan, os.Interrupt, os.Kill)
    return
}

func (sh *Handler)Bind(s os.Signal, cb Callback) {
    sh.cb[s] = cb
}

func (sh *Handler)Loop() os.Signal {
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
