package workerpool

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	testCases := []struct {
		Name        string
		Wokers      int
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name:        "ok",
			Wokers:      2,
			ErrorAssert: assert.NoError,
		},
		{
			Name:        "zero",
			Wokers:      0,
			ErrorAssert: assert.Error,
		},
		{
			Name:        "less than zero",
			Wokers:      -1,
			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			_, err := New(tc.Wokers)
			tc.ErrorAssert(t, err)
		})
	}
}

func Test_InitStop(t *testing.T) {
	wp, err := New(2)
	require.NoError(t, err)
	wp.Init()
	wp.Stop()
}

func Test_StopWaitTasks(t *testing.T) {
	wp, err := New(1)
	require.NoError(t, err)
	wp.Init()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	finished := make(chan struct{})
	start := make(chan struct{})
	err = wp.Add(ctx, func(_ context.Context) {
		<-start
		time.Sleep(50 * time.Millisecond)
		close(finished)
	})
	require.NoError(t, err)
	close(start)
	wp.Stop()
	select {
	case <-finished:
		return
	default:
		t.Fatalf("close should wait when task is finished")
	}
}

func Test_StopCancelTasks(t *testing.T) {
	wp, err := New(1)
	require.NoError(t, err)
	wp.Init()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	finished := make(chan struct{})
	start := make(chan struct{})
	err = wp.Add(ctx, func(fctx context.Context) {
		<-start
		<-fctx.Done()
		close(finished)
	})
	require.NoError(t, err)
	close(start)
	wp.Stop()
	select {
	case <-finished:
		return
	default:
		t.Fatalf("close should cancel running tasks")
	}
}

func Test_AddAfterStop(t *testing.T) {
	wp, err := New(1)
	require.NoError(t, err)
	wp.Init()
	wp.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = wp.Add(ctx, func(_ context.Context) {
		t.Fatalf("task should not be run if workerpool closed")
	})
	require.Error(t, err)
}

func Test_TasksAreConcurrent(t *testing.T) {
	wp, err := New(2)
	require.NoError(t, err)
	wp.Init()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	pipe := make(chan struct{})
	finished := make(chan struct{})
	err = wp.Add(ctx, func(_ context.Context) {
		pipe <- struct{}{}
	})
	require.NoError(t, err)
	err = wp.Add(ctx, func(_ context.Context) {
		<-pipe
		close(finished)
	})
	require.NoError(t, err)
	wp.Stop()
}

func Test_AddAlsoRespectContext(t *testing.T) {
	// only one place for task
	wp, err := New(1)
	require.NoError(t, err)
	wp.Init()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = wp.Add(ctx, func(fctx context.Context) {
		<-fctx.Done()
	})
	require.NoError(t, err)

	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2()
	err = wp.Add(ctx2, func(_ context.Context) {})
	require.Error(t, err)
	wp.Stop()
}
