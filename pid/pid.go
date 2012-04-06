package pid

import (
    "os"
    "strconv"
)

type PidFile struct {
    path string
    Pid int
}

func New(path string) (pf *PidFile, err error) {
    pf = &PidFile{path, os.Getpid()}
    f, err := os.OpenFile(pf.path, os.O_CREATE | os.O_EXCL | os.O_WRONLY, 0600)
    if err != nil {
        return
    }
    defer f.Close()
    _, err = f.WriteString(strconv.Itoa(pf.Pid))
    return
}

func (pf *PidFile)Close() error {
    return os.Remove(pf.path)
}

func (pf *PidFile)String() string {
    return strconv.Itoa(pf.Pid)
}
