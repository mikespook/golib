package autoinc

type AutoInc struct {
    start, step int
    queue chan int
    running bool
}

func New(start, step int) (ai *AutoInc) {
    ai = &AutoInc{
        start: start,
        step: step,
        running: true,
        queue: make(chan int, 4),
    }
    go ai.process()
    return
}

func (ai *AutoInc) process() {
    defer func() {recover()}()
    for i := ai.start; ai.running ; i=i+ai.step {
        ai.queue <- i
    }
}

func (ai *AutoInc) Id() int {
    return <-ai.queue
}

func (ai *AutoInc) Close() {
    ai.running = false
    close(ai.queue)
}
