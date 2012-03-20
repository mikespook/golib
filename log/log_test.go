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
    l, err := newLog("testing.log", LogAll)

    if err != nil {
        t.Error(err)
    }

    l.Error(errors.New("Test Error."))
    l.Warning("Test Warning.")
    l.Message("Test Message.")
    l.Debug("Test Debug.")
}

func TestDisableLog(t *testing.T) {
    defer os.Remove("testing.log")
    l, err := newLog("testing.log", LogAll ^ DisableDebug ^ DisableWarning)
    if err != nil {
        t.Error(err)
    }
    l.SetPrefix("{D}")

    l.Error(errors.New("Test Error."))
    l.Warning("Test Warning.")
    l.Message("Test Message.")
    l.Debug("Test Debug.")
}
