package scheduler

import (
    "time"
    "sync"
    "errors"
)

var (
    PoolSizePerTick = 10
    TickInterval = time.Second
    ErrTaskNotFound = errors.New("The task was not found.")
)

// Task Interface
type Task interface {
    // Start time
    Start() time.Duration
    SetStart(time.Duration)
    // interval of test executing, effective with Iterate() returns none-zero
    Interval() time.Duration
    // repeating times, 0 means don't repeat
    Iterate() int
    Id() interface{}
    Exec() error
    Cancel() error
}

type Scheduler struct {
    mutex sync.RWMutex
    ticks map[time.Duration][]interface{}
    tasks map[interface{}]Task
    HandleError func(error)
}

func New() *Scheduler {
    ts := &Scheduler {
        ticks: make(map[time.Duration][]interface{}),
        tasks: make(map[interface{}]Task),
    }
    return ts
}

func (ts *Scheduler) Put(task Task) {
    ts.mutex.Lock()
    defer ts.mutex.Unlock()
    id := task.Id()
    start := task.Start()
    ts.tasks[id] = task
    if ts.ticks[start] == nil {
        ts.ticks[start] = make([]interface{}, 0, PoolSizePerTick)
    }
    ts.ticks[start] = append(ts.ticks[start], id)
}

func (ts *Scheduler) Remove(id interface{}) {
    ts.mutex.Lock()
    defer ts.mutex.Unlock()
    delete(ts.tasks, id)
}

func (ts *Scheduler) Cancel(id interface{}) (err error) {
    ts.mutex.Lock()
    defer ts.mutex.Unlock()
    err = ts.tasks[id].Cancel()
    delete(ts.tasks, id)
    return
}

func (ts *Scheduler) Exec(id interface{}) (err error) {
    ts.mutex.Lock()
    defer ts.mutex.Unlock()
    if task, ok := ts.tasks[id]; ok {
        err = ts.exec(task)
    } else {
        err = ErrTaskNotFound
    }
    return
}

func (ts *Scheduler) Get(id interface{}) Task {
    ts.mutex.RLock()
    defer ts.mutex.RUnlock()
    return ts.tasks[id]
}

func (ts *Scheduler) Count() int {
    ts.mutex.RLock()
    defer ts.mutex.RUnlock()
    return len(ts.tasks)
}

func (ts *Scheduler) TickCount(t time.Duration) int {
    ts.mutex.RLock()
    defer ts.mutex.RUnlock()
    return len(ts.ticks[t])
}

func (ts *Scheduler) exec(task Task) (err error) {
    ts.mutex.Lock()
    defer func () {
        ts.mutex.Unlock()
        if e := recover(); e != nil {
            err = e.(error)
        }
    }()
    err = task.Exec()
    if task.Iterate() == 0 {
        delete(ts.tasks, task.Id())
    } else {
        task.SetStart(task.Start() + task.Interval())
        ts.Put(task)
    }
    return
}

func (ts *Scheduler) Loop() {
    for now := range time.Tick(TickInterval) {
        current := time.Duration(now.UnixNano())
        for t := range ts.ticks {
            // task executing time less/equal current time
            if t <= current {
                for index := range ts.ticks[t] {
                    id := ts.ticks[t][index]
                    if task, ok := ts.tasks[id]; ok {
                        go ts.exec(task)
                    }
                }
                delete(ts.ticks, t)
            }
        }
    }
}
