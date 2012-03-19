package util

import (
    "os"
    "strconv"
)

type PidFile string

func NewPidFile(path string) (pf PidFile, err error) {
    pf = PidFile(path)
    f, err := os.OpenFile(pf.String(), os.O_CREATE | os.O_EXCL | os.O_WRONLY, 0600)
    if err != nil {
        return
    }
    defer f.Close()
    _, err = f.WriteString(strconv.Itoa(os.Getpid()))
    return
}

func (pf PidFile)Close() error {
    return os.Remove(pf.String())
}

func (pf PidFile)String() string {
    return string(pf)
}

