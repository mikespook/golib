package util

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
    LogNone = 0
)

type Log struct {
    *log.Logger
    flag int
}

func NewLog(file string, flag int) (l *Log, err error){
    if file != "" {
        f, err := os.OpenFile(file, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600)
        if err == nil {
            l = &Log{log.New(f, "", log.LstdFlags), flag}
        }
    }
    if l == nil {
        l = &Log{log.New(os.Stdout, "", log.LstdFlags), flag}
    }
    return l, err
}

func (l *Log) Error(err error) {
    if l.flag ^ DisableError == 0 {
        return
    }
    l.Printf("[ERR] %s", err.Error())
}

func (l *Log) Warning(msg string) {
    if l.flag & DisableWarning == 0 {
        return
    }
    l.Printf("[WRN] %s", msg)
}

func (l *Log) Message(msg string) {
    if l.flag & DisableMessage == 0 {
        return
    }
    l.Printf("[MSG] %s", msg)
}

func (l *Log) Debug(msg string) {
    if l.flag & DisableDebug == 0 {
        return
    }
    l.Printf("[DBG] %s", msg)
}
