package scheduler

import (
    "time"
    "testing"
    "github.com/mikespook/golib/log"
    "github.com/mikespook/golib/autoinc"
)

var (
    ai = autoinc.New(0, 1)
)

type _task struct {
    start time.Duration
    iterate int
    interval time.Duration
    id int
}

func (t *_task)Start() time.Duration {
    return t.start
}

func (t *_task)SetStart(tm time.Duration) {
    t.start = tm
}

func (t *_task)Interval() time.Duration {
    return 0
}
func (t *_task)Iterate() int {
    return 0
}

func (t *_task)Id() interface{} {
    return t.id
}
func (t *_task)Exec() error {
    log.Debugf("Task %d Executed.", t.Id())
    return nil
}
func (t *_task)Cancel() error {
    return nil
}

func Test(t *testing.T) {
    ts := New()
    ts.HandleError = func(err error) {
        t.Error(err)
    }
    go ts.Loop()
    n := time.Duration(time.Now().UnixNano())
    ts.Put(&_task{
        id: ai.Id(),
        start: n + time.Second,
        interval: time.Second,
        iterate: 0,
    })
    ts.Put(&_task{
        id: ai.Id(),
        start: n + time.Second,
        interval: time.Second,
        iterate: 0,
    })
    ts.Put(&_task{
        id: ai.Id(),
        start: n + 3 * time.Second,
        interval: time.Second,
        iterate: 0,
    })
    ts.Put(&_task{
        id: ai.Id(),
        start: n + 3 * time.Second,
        interval: time.Second,
        iterate: 0,
    })
    c := ts.Count()
    if c != 4 {
        t.Errorf("Task count should be 4 but get %d.", c)
    }
    c = ts.TickCount(n + time.Second)
    if c != 2 {
        t.Errorf("Task count should be 2 but get %d.", c)
    }
    c = ts.TickCount(n + 3 * time.Second)
    if c != 2 {
        t.Errorf("Task count should be 2 but get %d.", c)
    }
    time.Sleep(2 * time.Second)
    c = ts.Count()
    if c != 2 {
        t.Errorf("Task count should be 2 but get %d.", c)
    }
    c = ts.TickCount(n)
    if c != 0 {
        t.Errorf("Task count should be 0 but get %d.", c)
    }
    c = ts.TickCount(n + 3 * time.Second)
    if c != 2 {
        t.Errorf("Task count should be 2 but get %d.", c)
    }
    time.Sleep(2 * time.Second)
    c = ts.Count()
    if c != 0 {
        t.Errorf("Task count should be 0 but get %d.", c)
    }
    c = ts.TickCount(n)
    if c != 0 {
        t.Errorf("Task count should be 0 but get %d.", c)
    }
    c = ts.TickCount(n + 2 * time.Second)
    if c != 0 {
        t.Errorf("Task count should be 0 but get %d.", c)
    }

}
