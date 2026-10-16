// Command worker consumes tasks and sends them to an LLM.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sachonidas/llm-task-worker/internal/llm/llmtest"
	"github.com/sachonidas/llm-task-worker/internal/task"
	"github.com/sachonidas/llm-task-worker/internal/worker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tasks := make(chan task.Task)

	go func() {
		defer close(tasks) // quien escribe, cierra
		for i := 1; i <= 10; i++ {
			t := task.Task{ID: strconv.Itoa(i), Prompt: "hi"}
			select {
			case tasks <- t:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Consumidor: el pool.
	results := worker.Run(ctx, 3, tasks, &llmtest.Fake{Delay: time.Second})

	for r := range results {
		if r.Err != nil {
			log.Printf("task %s: error: %v", r.Task.ID, r.Err)
			continue
		}
		log.Printf("task %s: %s", r.Task.ID, r.Response.Text)
	}

	log.Println("shutdown complete")
}
