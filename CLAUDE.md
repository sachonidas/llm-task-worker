# llm-task-worker

Go worker that consumes tasks from a queue and calls an LLM API reliably.
Goal: I'm a senior PHP dev learning idiomatic Go. Portfolio project.

## Rules

- Teaching mode: explain Go idioms before writing code, and compare with PHP (Fibers, exceptions, DI) when useful.
- Core concurrency code (worker pool, channels, context, shutdown): guide me with hints, let me write it, then review it.
- Scaffolding, boilerplate, fakes and docker-compose: you can write them, but explain what you did.
- Idiomatic Go only: small interfaces, errors as values wrapped with %w, standard library first, no frameworks.
- Tests for everything; use fakes via interfaces, not mock libraries.
- Small commits with clear messages.

## Phases

1. Worker pool with an in-memory channel: task type, LLM client interface + fake, pool, graceful shutdown with context and signal.NotifyContext.
2. RabbitMQ: consume tasks from a queue (docker-compose), ack/nack correctly.
3. Retries with exponential backoff + dead-letter queue.
4. Idempotency: never process the same task twice.
5. Rate limiting and per-call cost tracking.
6. Observability: structured logging (slog) + Prometheus metrics.
7. Real Claude API client.
