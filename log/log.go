package log

import (
    "os"
    "log"
    "fmt"
    "errors"
)

const (
    DisableError = 1
    DisableWarning = 2
    DisableMessage = 4
    DisableDebug = 8
    LogAll = 0xF
    LogNone = 0
    LogError = LogAll ^ DisableWarning ^ DisableMessage ^ DisableDebug
    LogWarning = LogAll ^ DisableMessage ^ DisableDebug ^ DisableError
    LogMessage = LogAll ^ DisableDebug ^ DisableError ^ DisableWarning
    LogDebug = LogAll ^ DisableError ^ DisableWarning ^ DisableMessage
)

var (
    DefaultLogger *Logger
)

func init() {
    DefaultLogger, _ = newLog("", LogAll)
}

func Init(file string, flag int) (err error) {
    DefaultLogger, err = newLog(file, flag)
    return
}

type Logger struct {
    *log.Logger
    flag int
}

func newLog(file string, flag int) (l *Logger, err error){
    if file != "" {
        f, err := os.OpenFile(file, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600)
        if err == nil {
            l = &Logger{log.New(f, "", log.LstdFlags), flag}
        }
    }
    if l == nil {
        l = &Logger{log.New(os.Stdout, "", log.LstdFlags), flag}
    }
    return l, err
}

func Errorf(format string, msg ... interface{}) {
    DefaultLogger.Errorf(format, msg ... )
}

func (l *Logger) Errorf(format string, msg ... interface{}) {
    l.Error(errors.New(fmt.Sprintf(format, msg ...)))
}

func Error(err error) {
    DefaultLogger.Error(err)
}

func (l *Logger) Error(err error) {
    if l.flag & DisableError == 0 {
        return
    }
    l.Printf("[ERR] %s", err.Error())
}

func Warning(msg ... interface{}) {
    DefaultLogger.Warning(msg ... )
}

func (l *Logger) Warning(msg ... interface{}) {
    if l.flag & DisableWarning == 0 {
        return
    }
    l.Printf("[WRN] %s", msg ...)
}

func Warningf(format string, msg ... interface{}) {
    DefaultLogger.Warningf(format, msg ... )
}

func (l *Logger) Warningf(format string, msg ... interface{}) {
    l.Warning(fmt.Sprintf(format, msg ...))
}

func Message(msg ... interface{}) {
    DefaultLogger.Message(msg ... )
}

func (l *Logger) Message(msg ... interface{}) {
    if l.flag & DisableMessage == 0 {
        return
    }
    l.Printf("[MSG] %s", msg ...)
}

func Messagef(format string, msg ... interface{}) {
    DefaultLogger.Messagef(format, msg ... )
}


func (l *Logger) Messagef(format string, msg ... interface{}) {
    l.Message(fmt.Sprintf(format, msg ...))
}

func Debug(msg ... interface{}) {
    DefaultLogger.Debug(msg ... )
}

func (l *Logger) Debug(msg ... interface{}) {
    if l.flag & DisableDebug == 0 {
        return
    }
    l.Printf("[DBG] %s", msg ...)
}

func Debugf(format string, msg ... interface{}) {
    DefaultLogger.Debugf(format, msg ... )
}

func (l *Logger) Debugf(format string, msg ... interface{}) {
    l.Debug(fmt.Sprintf(format, msg ...))
}
