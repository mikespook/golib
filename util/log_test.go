package util

import (
    "errors"
    "testing"
)

func TestFlag(t *testing.T) {
    f1 := LogAll ^ DisableDebug
    if f1 & DisableDebug != 0 {
        t.Error("Flags of log was broken.")
    }
}

func TestNewLog(t *testing.T) {
    var err error

    _, err = NewLog("", LogAll)
    if err != nil {
        t.Error(err)
    }

    _, err = NewLog("testing.log", LogAll)
    if err != nil {
        t.Error(err)
    }

    _, err = NewLog("foobar/testing.log", LogAll)
    if err != nil {
        t.Log(err)
    }
}

func TestLog(t *testing.T) {
    l, err := NewLog("testing.log", LogAll)

    if err != nil {
        t.Error(err)
    }

    l.Error(errors.New("Test Error."))
    l.Warning("Test Warning.")
    l.Message("Test Message.")
    l.Debug("Test Debug.")
}

func TestDisableLog(t *testing.T) {
    l, err := NewLog("testing.log", LogAll ^ DisableDebug ^ DisableWarning)
    if err != nil {
        t.Error(err)
    }
    l.SetPrefix("{D}")

    l.Error(errors.New("Test Error."))
    l.Warning("Test Warning.")
    l.Message("Test Message.")
    l.Debug("Test Debug.")
}
