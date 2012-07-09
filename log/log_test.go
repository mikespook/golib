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

    a, err := NewLog("", LogAll, DefaultBufSize)
    if err != nil {
        t.Error(err)
    }
    defer a.Close()

    b, err := NewLog("testing.log", LogAll, DefaultBufSize)
    if err != nil {
        t.Error(err)
    }
    defer b.Close()

    c, err := NewLog("foobar/testing.log", LogAll, DefaultBufSize)
    if err != nil {
        t.Log(err)
    }
    defer c.Close()
}

func TestLog(t *testing.T) {
    defer os.Remove("testing.log")
    err := Init("testing.log", LogAll)
    if err != nil {
        t.Error(err)
    }
    defer Close()

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
    defer Close()

    Error(errors.New("Test Error."))
    Warning("Test Warning.")
    Message("Test Message.")
    Debug("Test Debug.")
}

