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
    var f *os.File
    f, err = os.OpenFile(path, os.O_CREATE | os.O_EXCL | os.O_WRONLY, 0600)
    if err != nil {
        f, err = os.OpenFile(path, os.O_CREATE | os.O_RDWR, 0600)
        if err != nil {
            return
        }
        var pidstr [10]byte
        var n int
        n, err = f.Read(pidstr[:])
        if err != nil {
            return
        }
        var pid int
        pid, err = strconv.Atoi(string(pidstr[:n]))
        if err != nil {
            return
        }
        _, err = os.FindProcess(pid)
        if err == nil {
            // process exists
            err = ProcessExists
            return
        }
        f.Truncate(int64(n))
    }
    defer f.Close()
    _, err = f.WriteString(strconv.Itoa(pf.Pid))
    return
}

func (pf *PidFile)Close() error {
    return os.Remove(pf.path)
}
