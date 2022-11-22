package workerpool

import (
	"context"
	"errors"
	"sync"
)

// WorkerPool is simple worker pool with preallocated goroutines for tasks
// and tasks cancelation on workerpool close.
type WorkerPool struct {
	workersCount int
	tasks        chan func()

	stopOnce     sync.Once
	stopSignal   chan struct{}
	wg           sync.WaitGroup
	mu           sync.Mutex
	stopped      bool
	currentTasks map[context.Context]context.CancelFunc
}

func New(workers int) (*WorkerPool, error) {
	if workers <= 0 {
		return nil, errors.New("number of workers should be greater than zero")
	}
	return &WorkerPool{
		workersCount: workers,
		tasks:        make(chan func()),
		stopSignal:   make(chan struct{}),
		currentTasks: map[context.Context]context.CancelFunc{},
	}, nil
}

func (wp *WorkerPool) Init() {
	for i := 0; i < wp.workersCount; i++ {
		wp.wg.Add(1)
		go wp.runTasks()
	}
}

func (wp *WorkerPool) Stop() {
	wp.stopOnce.Do(func() {
		close(wp.stopSignal)
		wp.mu.Lock()
		wp.stopped = true
		for _, cancel := range wp.currentTasks {
			cancel()
		}
		wp.mu.Unlock()
	})
	wp.wg.Wait()
}

func (wp *WorkerPool) runTasks() {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.stopSignal:
			return
		case task := <-wp.tasks:
			task()
		}
	}
}

func (wp *WorkerPool) Add(ctx context.Context, fn func(context.Context)) error {
	ctx, cancel := context.WithCancel(ctx)
	task := func() {
		fn(ctx)
		cancel()
	}
	wp.mu.Lock()
	if wp.stopped {
		defer wp.mu.Unlock()
		return errors.New("worker pool is closed")
	}
	wp.currentTasks[ctx] = cancel
	wp.mu.Unlock()
	// we don't need to check stopSignal here, because if workerpool is closed after
	// we reached mutex then ctx will be cancelled anyway and we can safely return from Add.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case wp.tasks <- task:
	}
	return nil
}
