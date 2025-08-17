package internal

import (
	"sync/atomic"
	"time"
)

type EventLoopTask func()

type EventLoop interface {
	SetPostTimeout(timeout time.Duration)
	PostTask(task EventLoopTask)
	Run()
	RunUntilIdle()
	Quit()
}

type eventLoopImpl struct {
	tasks   chan EventLoopTask
	timeout time.Duration
	running atomic.Bool
}

func CreateEventLoop(size int) *eventLoopImpl {
	result := eventLoopImpl{
		tasks: make(chan EventLoopTask, size),
	}
	result.running.Store(true)
	result.timeout = time.Second * 10
	return &result
}

func (el *eventLoopImpl) SetPostTimeout(timeout time.Duration) {
	el.timeout = timeout
}

func (el *eventLoopImpl) PostTask(task EventLoopTask) {
	select {
	case el.tasks <- task:
	case <-time.After(el.timeout):
		panic("Event loop was become irresponsible")
	}
}

func (el *eventLoopImpl) Run() {
	for el.running.Load() {
		task := <-el.tasks
		task()
	}
}

func (el *eventLoopImpl) RunUntilIdle() {
	for {
		select {
		case task := <-el.tasks:
			task()
		default:
			return
		}
	}
}

func (el *eventLoopImpl) Quit() {
	el.running.Store(false)
	el.tasks <- func() {}
}
