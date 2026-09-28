// Package llmtest provides test doubles for llm.Client, in the spirit of
// net/http/httptest.
package llmtest

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/sachonidas/llm-task-worker/internal/llm"
)

// Compile-time check that *Fake implements llm.Client.
var _ llm.Client = (*Fake)(nil)

// Fake is a configurable llm.Client. The zero value is ready to use and
// echoes the prompt immediately. Configure fields before sharing it
// between goroutines.
type Fake struct {
	// Delay simulates network latency. Cancelling ctx aborts the wait.
	Delay time.Duration
	// Err, if set, is returned by every call.
	Err error

	calls atomic.Int64
}

// Complete implements llm.Client.
func (f *Fake) Complete(ctx context.Context, req llm.Request) (llm.Response, error) {
	f.calls.Add(1)

	select {
	case <-ctx.Done():
		return llm.Response{}, ctx.Err()
	case <-time.After(f.Delay):
	}

	if f.Err != nil {
		return llm.Response{}, f.Err
	}
	return llm.Response{
		Text:         "echo: " + req.Prompt,
		InputTokens:  len(req.Prompt),
		OutputTokens: len(req.Prompt) + len("echo: "),
	}, nil
}

// Calls returns how many times Complete was called.
func (f *Fake) Calls() int {
	return int(f.calls.Load())
}
