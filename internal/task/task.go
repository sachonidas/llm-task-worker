// Package task defines the unit of work the worker processes.
package task

import "errors"

// Sentinel errors: callers compare with errors.Is, even if wrapped with %w.
var (
	ErrEmptyID     = errors.New("task: empty id")
	ErrEmptyPrompt = errors.New("task: empty prompt")
)

// Task is a single job to send to the LLM.
type Task struct {
	ID     string
	Prompt string
}

// Validate reports whether the task can be processed.
func (t Task) Validate() error {
	if t.ID == "" {
		return ErrEmptyID
	}
	if t.Prompt == "" {
		return ErrEmptyPrompt
	}
	return nil
}
