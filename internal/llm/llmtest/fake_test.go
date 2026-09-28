package llmtest

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/sachonidas/llm-task-worker/internal/llm"
)

func TestFake_Echo(t *testing.T) {
	f := &Fake{}
	resp, err := f.Complete(context.Background(), llm.Request{Prompt: "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Text != "echo: hi" {
		t.Errorf("Text = %q, want %q", resp.Text, "echo: hi")
	}
}

func TestFake_Err(t *testing.T) {
	boom := errors.New("boom")
	f := &Fake{Err: boom}
	_, err := f.Complete(context.Background(), llm.Request{Prompt: "hi"})
	if !errors.Is(err, boom) {
		t.Errorf("err = %v, want %v", err, boom)
	}
}

func TestFake_CancelDuringDelay(t *testing.T) {
	f := &Fake{Delay: time.Hour}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := f.Complete(ctx, llm.Request{Prompt: "hi"})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want %v", err, context.DeadlineExceeded)
	}
}

func TestFake_ConcurrentCalls(t *testing.T) {
	f := &Fake{}
	const n = 100

	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			f.Complete(context.Background(), llm.Request{Prompt: "hi"})
		})
	}
	wg.Wait()

	if got := f.Calls(); got != n {
		t.Errorf("Calls() = %d, want %d", got, n)
	}
}
