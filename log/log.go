package log

import (
    "os"
    "log"
)

const (
    DisableError = 1
    DisableWarning = 2
    DisableMessage = 4
    DisableDebug = 8
    LogAll = 0xF
    None = 0
    LogError = LogAll ^ DisableWarning ^ DisableMessage ^ DisableDebug
    LogWarning = LogAll ^ DisableMessage ^ DisableDebug ^ DisableError
    LogMessage = LogAll ^ DisableDebug ^ DisableError ^ DisableWarning
    LogDebug = LogAll ^ DisableError ^ DisableWarning ^ DisableMessage
)

var (
    l *logger
)

func Init(file string, flag int) (err error) {
    l, err = newLog(file, flag)
    return
}

type logger struct {
    *log.Logger
    flag int
}

func newLog(file string, flag int) (l *logger, err error){
    if file != "" {
        f, err := os.OpenFile(file, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600)
        if err == nil {
            l = &logger{log.New(f, "", log.LstdFlags), flag}
        }
    }
    if l == nil {
        l = &logger{log.New(os.Stdout, "", log.LstdFlags), flag}
    }
    return l, err
}

func Error(err error) {
    if l.flag ^ DisableError == 0 {
        return
    }
    l.Printf("[ERR] %s", err.Error())
}

func Warning(msg string) {
    if l.flag & DisableWarning == 0 {
        return
    }
    l.Printf("[WRN] %s", msg)
}

func Message(msg string) {
    if l.flag & DisableMessage == 0 {
        return
    }
    l.Printf("[MSG] %s", msg)
}

func Debug(msg string) {
    if l.flag & DisableDebug == 0 {
        return
    }
    l.Printf("[DBG] %s", msg)
}
