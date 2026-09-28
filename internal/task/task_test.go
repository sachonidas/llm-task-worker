package task

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want error
	}{
		{"valid", Task{ID: "1", Prompt: "hi"}, nil},
		{"empty id", Task{Prompt: "hi"}, ErrEmptyID},
		{"empty prompt", Task{ID: "1"}, ErrEmptyPrompt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.task.Validate(); !errors.Is(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
