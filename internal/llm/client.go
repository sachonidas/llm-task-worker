// Package llm defines the contract for talking to a language model.
package llm

import "context"

// Request is the input for a completion. It will grow (model, max tokens...)
// without breaking the Client signature.
type Request struct {
	Prompt string
}

// Response is the model output. Token counts feed cost tracking later.
type Response struct {
	Text         string
	InputTokens  int
	OutputTokens int
}

// Client sends a prompt to an LLM. Implementations must honor ctx
// cancellation and be safe for concurrent use.
type Client interface {
	Complete(ctx context.Context, req Request) (Response, error)
}
