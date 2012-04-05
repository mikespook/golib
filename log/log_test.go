package log

import (
    "errors"
    "testing"
    "os"
)

func TestFlag(t *testing.T) {
    f1 := LogAll ^ DisableDebug
    if f1 & DisableDebug != 0 {
        t.Error("Flags of log was broken.")
    }
}

func TestNewLog(t *testing.T) {
    defer os.Remove("testing.log")
    var err error

    _, err = newLog("", LogAll)
    if err != nil {
        t.Error(err)
    }

    _, err = newLog("testing.log", LogAll)
    if err != nil {
        t.Error(err)
    }

    _, err = newLog("foobar/testing.log", LogAll)
    if err != nil {
        t.Log(err)
    }
}

func TestLog(t *testing.T) {
    defer os.Remove("testing.log")
    err := Init("testing.log", LogAll)
    if err != nil {
        t.Error(err)
    }

    Error(errors.New("Test Error."))
    Warning("Test Warning.")
    Message("Test Message.")
    Debug("Test Debug.")
}

func TestDisableLog(t *testing.T) {
    defer os.Remove("testing.log")
    err := Init("testing.log", LogAll ^ DisableDebug ^ DisableWarning)
    if err != nil {
        t.Error(err)
    }

    Error(errors.New("Test Error."))
    Warning("Test Warning.")
    Message("Test Message.")
    Debug("Test Debug.")
}
