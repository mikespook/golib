package pid

import (
    "os"
    "errors"
    "strconv"
)

var (
    ProcessExists = errors.New("Process already exists.")
)

type PidFile struct {
    path string
    Pid int
}

func New(path string) (pf *PidFile, err error) {
    pf = &PidFile{path, os.Getpid()}
    oldf, err := os.OpenFile(pf.path, os.O_CREATE | os.O_EXCL | os.O_WRONLY, 0600)
    if err != nil {
        var pidstr [10]byte
        n, err := oldf.Read(pidstr[:])
        if err != nil {
            return nil, err
        }
        pid, err := strconv.Atoi(string(pidstr[:n]))
        if err != nil {
            return nil, err
        }
        _, err = os.FindProcess(pid)
        if err == nil {
            // process exists
            return nil, ProcessExists
        }
    }
    f, err := os.OpenFile(pf.path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0600)
    if err != nil {
        return nil, err
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
