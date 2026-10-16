package worker

import (
	"context"
	"sync"

	"github.com/sachonidas/llm-task-worker/internal/llm"
	"github.com/sachonidas/llm-task-worker/internal/task"
)

type Result struct {
	Task     task.Task
	Response llm.Response
	Err      error
}

func Run(ctx context.Context, workers int, tasks <-chan task.Task, client llm.Client) <-chan Result {
	results := make(chan Result, workers)

	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-tasks:
					if !ok {
						return
					}
					response, err := client.Complete(ctx, llm.Request{Prompt: t.Prompt})
					r := Result{
						Task:     t,
						Response: response,
						Err:      err,
					}

					select {
					case results <- r:
					case <-ctx.Done():
						return
					}
				}
			}
		})
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
