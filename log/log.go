package log

import (
    "os"
    "log"
    "fmt"
    "errors"
    "sync"
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
    DefaultBufSize = 32
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
    logChan chan *logRecord
    mutex sync.Mutex
    f *os.File
}

type logRecord struct {
    Type uint8
    Message []interface{}
}

func NewLog(file string, flag, bufsize int) (l *Logger, err error){
    if file != "" {
        f, err := os.OpenFile(file, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600)
        if err == nil {
            l = &Logger{Logger:log.New(f, "", log.LstdFlags), flag:flag, f: f}
        }
    }
    if l == nil {
        l = &Logger{Logger:log.New(os.Stdout, "", log.LstdFlags), flag:flag, f: os.Stdout}
    }
    l.logChan = make(chan *logRecord, bufsize)
    go func() {
        l.mutex.Lock()
        defer l.mutex.Unlock()
        var t string
        for record := range l.logChan {
            switch record.Type {
                case TypeDebug:
                    t = "[DBG] %s"
                case TypeMessage:
                    t = "[MSG] %s"
                case TypeWarning:
                    t = "[WRN] %s"
                case TypeError:
                    t = "[ERR] %s"
            }
            l.Printf(t, record.Message ... )
        }
    } ()
    return l, err
}

func(l *Logger) Fd() uintptr {
    return l.f.Fd()
}

func (l *Logger) Errorf(format string, msg ... interface{}) {
    l.Error(errors.New(fmt.Sprintf(format, msg ...)))
}

func (l *Logger) Error(err error) {
    if l.flag & DisableError == 0 {
        return
    }
    l.logChan <- &logRecord{Type:TypeError, Message: []interface{}{err}}
}

func (l *Logger) Warning(msg ... interface{}) {
    if l.flag & DisableWarning == 0 {
        return
    }
    l.logChan <- &logRecord{Type:TypeWarning, Message: msg}
}

func (l *Logger) Warningf(format string, msg ... interface{}) {
    l.Warning(fmt.Sprintf(format, msg ...))
}

func (l *Logger) Message(msg ... interface{}) {
    if l.flag & DisableMessage == 0 {
        return
    }
    l.logChan <- &logRecord{Type:TypeMessage, Message: msg}
}

func (l *Logger) Messagef(format string, msg ... interface{}) {
    l.Message(fmt.Sprintf(format, msg ...))
}

func (l *Logger) Debug(msg ... interface{}) {
    if l.flag & DisableDebug == 0 {
        return
    }
    l.logChan <- &logRecord{Type:TypeDebug, Message: msg}
}

func (l *Logger) Debugf(format string, msg ... interface{}) {
    l.Debug(fmt.Sprintf(format, msg ...))
}

func (l *Logger) Close() {
    close(l.logChan)
}

// Close the logger and waiting all messages was printed
func (l *Logger) WaitClosing() {
    defer l.mutex.Unlock()
    l.Close()
    l.mutex.Lock()
}

var (
    DefaultLogger *Logger
)

func init() {
    DefaultLogger, _ = NewLog("", LogAll, DefaultBufSize)
}

func Fd() uintptr {
    return DefaultLogger.Fd()
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

func Close() {
    DefaultLogger.Close()
}

func WaitClosing() {
    DefaultLogger.WaitClosing()
}

func Exit(code int) {
    runtime.Gosched()
    os.Exit(code)
}
