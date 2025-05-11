package internal

import "sync/atomic"

type EventLoopTask func()

type EventLoop interface {
	PostTask(task EventLoopTask)
	Run()
	Quit()
}

type eventLoopImpl struct {
	tasks   chan EventLoopTask
	running atomic.Bool
}

func CreateEventLoop() *eventLoopImpl {
	result := eventLoopImpl{
		tasks: make(chan EventLoopTask, 40),
	}
	result.running.Store(true)
	return &result
}

func (el *eventLoopImpl) PostTask(task EventLoopTask) {
	el.tasks <- task
}

func (el *eventLoopImpl) Run() {
	for el.running.Load() {
		task := <-el.tasks
		task()
	}
}

func (el *eventLoopImpl) Quit() {
	el.running.Store(false)
	el.tasks <- func() {}
}
