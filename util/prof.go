package util

import (
    "os"
    "runtime/pprof"
)

func StartProfile(fn string) (err error) {
    f , err := os.Create(fn)
    if err != nil {
        return
    }
    pprof.StartCPUProfile(f)
    return
}

func StopProfile() {
    pprof.StopCPUProfile()
}

