package log

import (
    "io"
    "os"
    "log"
    "fmt"
    "errors"
    "runtime"
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

const (
    TypeDebug = iota
    TypeMessage
    TypeWarning
    TypeError
)

type Logger struct {
    *log.Logger
    flag int
}

func New(w io.Writer, flag, bufsize int) (l *Logger, err error) {
    l = &Logger{Logger:log.New(w, "", log.LstdFlags), flag:flag}
    return l, err
}

func NewLog(file string, flag, bufsize int) (l *Logger, err error){
    var f *os.File
    if file != "" {
        f, err = os.OpenFile(file, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600)
        if err != nil {
            f = os.Stdout
        }
    }
    if f == nil {
        f = os.Stdout
    }
    return New(f, flag, bufsize)
}

func (l *Logger) Errorf(format string, msg ... interface{}) {
    l.Error(errors.New(fmt.Sprintf(format, msg ...)))
}

func (l *Logger) Error(err error) {
    if l.flag & DisableError == 0 {
        return
    }
    l.Printf("[ERR] %s", err)
}

func (l *Logger) Warning(msg ... interface{}) {
    if l.flag & DisableWarning == 0 {
        return
    }
    l.Printf("[WRN] %s", msg ...)
}

func (l *Logger) Warningf(format string, msg ... interface{}) {
    l.Warning(fmt.Sprintf(format, msg ...))
}

func (l *Logger) Message(msg ... interface{}) {
    if l.flag & DisableMessage == 0 {
        return
    }
    l.Printf("[MSG] %s", msg ...)
}

func (l *Logger) Messagef(format string, msg ... interface{}) {
    l.Message(fmt.Sprintf(format, msg ...))
}

func (l *Logger) Debug(msg ... interface{}) {
    if l.flag & DisableDebug == 0 {
        return
    }
	l.Printf("[DBG] %s", msg ...)
}

func (l *Logger) Debugf(format string, msg ... interface{}) {
    l.Debug(fmt.Sprintf(format, msg ...))
}

var (
    DefaultLogger *Logger
    DefaultBufSize = 32
)

func init() {
    DefaultLogger, _ = NewLog("", LogAll, DefaultBufSize)
}

func Init(file string, flag int) (err error) {
    DefaultLogger, err = NewLog(file, flag, DefaultBufSize)
    return
}

func Error(err error) {
    DefaultLogger.Error(err)
}

func Errorf(format string, msg ... interface{}) {
    DefaultLogger.Errorf(format, msg ... )
}

func Warning(msg ... interface{}) {
    DefaultLogger.Warning(msg ... )
}

func Warningf(format string, msg ... interface{}) {
    DefaultLogger.Warningf(format, msg ... )
}

func Message(msg ... interface{}) {
    DefaultLogger.Message(msg ... )
}

func Messagef(format string, msg ... interface{}) {
    DefaultLogger.Messagef(format, msg ... )
}

func Debug(msg ... interface{}) {
    DefaultLogger.Debug(msg ... )
}

func Debugf(format string, msg ... interface{}) {
    DefaultLogger.Debugf(format, msg ... )
}

func Exit(code int) {
    runtime.Gosched()
    os.Exit(code)
}
