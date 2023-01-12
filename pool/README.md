# Go Pool

## Pool

Implementation

```go
package gopool

import (
	"fmt"
	"time"
)

// ErrScheduleTimeout returned by Pool to indicate that there no free
// goroutines during some period of time.
var ErrScheduleTimeout = fmt.Errorf("schedule error: timed out")

// Pool contains logic of goroutine reuse.
// semaphore, 用于控制协程的信号.
// 信号量(英语: semaphore) 又称为信号标, 是一个同步对象, 用于保持在0至指定最大值之间的一个计数值. 当线程完成一次对该semaphore对象的等待(wait)时, 该计数值减一; 当线程完成一次对semaphore对象的释放(release)时, 计数值加一.
type Pool struct {
	sem  chan struct{}
	work chan func()
}

// NewPool creates new goroutine pool with given size. It also creates a work
// queue of given size. Finally, it spawns given amount of goroutines
// immediately.
func NewPool(size, queue, spawn int) *Pool {
	if spawn <= 0 && queue > 0 {
		panic("dead queue configuration detected")
	}
	if spawn > size {
		panic("spawn > workers")
	}
	p := &Pool{
		sem:  make(chan struct{}, size),
		work: make(chan func(), queue),
	}
	for i := 0; i < spawn; i++ {
		p.sem <- struct{}{}
		go p.worker(func() {})
	}

	return p
}

// Schedule schedules task to be executed over pool's workers.
func (p *Pool) Schedule(task func()) {
	p.schedule(task, nil)
}

// ScheduleTimeout schedules task to be executed over pool's workers.
// It returns ErrScheduleTimeout when no free workers met during given timeout.
func (p *Pool) ScheduleTimeout(timeout time.Duration, task func()) error {
	return p.schedule(task, time.After(timeout))
}

func (p *Pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrScheduleTimeout
	case p.work <- task:
		return nil
	case p.sem <- struct{}{}:
		go p.worker(task)
		return nil
	}
}

func (p *Pool) worker(task func()) {
	defer func() { <-p.sem }()

	task()

	for task := range p.work {
		task()
	}
}
```

Traditional model

```go
go func(){
    doXXX()
}()
```

Use pool

```go
pool := gopool.NewPool(1024, 1024, 1024) // 限制1024个协程
pool.Schedule(func () {
    doXXX()
})
```

Benchmark result

```log
BenchmarkNaive
BenchmarkNaive-8   	   10000	      1136 ns/op	     225 B/op	       2 allocs/op
BenchmarkPool
BenchmarkPool-8    	   10000	       727.9 ns/op	     107 B/op	       2 allocs/op

BenchmarkNaive
BenchmarkNaive-8   	 2458435	       421.6 ns/op	      96 B/op	       2 allocs/op
BenchmarkPool
BenchmarkPool-8    	 1739314	       708.9 ns/op	      16 B/op	       1 allocs/op
```
