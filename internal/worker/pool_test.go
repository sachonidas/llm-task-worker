package worker

import (
	"context"
	"testing"
	"time"

	"github.com/sachonidas/llm-task-worker/internal/llm/llmtest"
	"github.com/sachonidas/llm-task-worker/internal/task"
)

func TestRun_ProcessesAllTasks(t *testing.T) {
	// Arrange
	tasks := make(chan task.Task)
	go func() {
		for _, id := range []string{"1", "2", "3"} {
			tasks <- task.Task{ID: id, Prompt: "hi"}
		}
		close(tasks) // quien escribe, cierra
	}()

	// Act
	results := Run(context.Background(), 2, tasks, &llmtest.Fake{})

	// Assert
	count := 0
	for r := range results {
		if r.Err != nil {
			t.Errorf("task %s: unexpected error: %v", r.Task.ID, r.Err)
		}
		count++
	}
	if count != 3 {
		t.Errorf("got %d results, want 3", count)
	}
}

func TestRun_StopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//create channel
	tasks := make(chan task.Task)

	results := Run(ctx, 2, tasks, &llmtest.Fake{Delay: time.Hour})

	go func() {
		tasks <- task.Task{ID: "test", Prompt: "hi"}
	}()

	cancel()

	done := make(chan struct{})
	go func() {
		for range results {
		}
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		t.Fatal("pool did not stop after context cancellation")
	}
}
