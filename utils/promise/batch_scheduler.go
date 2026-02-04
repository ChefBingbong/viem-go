// Package promise provides utilities for async operations and batching.
package promise

import (
	"context"
	"sync"
	"time"
)

// BatchSchedulerConfig contains configuration for the batch scheduler.
type BatchSchedulerConfig[P any, R any] struct {
	// ID is a unique identifier for this scheduler instance.
	ID string

	// Fn is the function that processes a batch of items.
	Fn func(ctx context.Context, args []P) ([]R, error)

	// ShouldSplitBatch determines if the current batch should be split.
	// Called before adding a new item to the batch.
	ShouldSplitBatch func(args []P) bool

	// Wait is the duration to wait before executing a batch.
	// Default: 0 (execute immediately when ready)
	Wait time.Duration
}

// BatchScheduler manages batched execution of requests.
type BatchScheduler[P any, R any] struct {
	config  BatchSchedulerConfig[P, R]
	mu      sync.Mutex
	pending []pendingItem[P, R]
	timer   *time.Timer
}

type pendingItem[P any, R any] struct {
	arg    P
	result chan batchResult[R]
}

type batchResult[R any] struct {
	value R
	err   error
}

// schedulerCache holds global scheduler instances by ID
var schedulerCache sync.Map

// CreateBatchScheduler creates or retrieves a batch scheduler with the given configuration.
// Schedulers are cached by ID for reuse.
func CreateBatchScheduler[P any, R any](config BatchSchedulerConfig[P, R]) *BatchScheduler[P, R] {
	// Check cache first
	if cached, ok := schedulerCache.Load(config.ID); ok {
		return cached.(*BatchScheduler[P, R])
	}

	scheduler := &BatchScheduler[P, R]{
		config:  config,
		pending: make([]pendingItem[P, R], 0),
	}

	// Store in cache
	schedulerCache.Store(config.ID, scheduler)

	return scheduler
}

// Schedule adds an item to the batch and returns when processed.
func (s *BatchScheduler[P, R]) Schedule(ctx context.Context, arg P) (R, error) {
	s.mu.Lock()

	// Check if we should split the batch before adding
	if s.config.ShouldSplitBatch != nil {
		args := make([]P, len(s.pending)+1)
		for i, p := range s.pending {
			args[i] = p.arg
		}
		args[len(s.pending)] = arg

		if len(s.pending) > 0 && s.config.ShouldSplitBatch(args) {
			// Execute current batch first
			s.executeNow()
		}
	}

	// Create result channel
	resultCh := make(chan batchResult[R], 1)
	item := pendingItem[P, R]{
		arg:    arg,
		result: resultCh,
	}

	// Add to pending
	wasEmpty := len(s.pending) == 0
	s.pending = append(s.pending, item)

	// Start timer if this is the first item
	if wasEmpty {
		if s.config.Wait > 0 {
			s.timer = time.AfterFunc(s.config.Wait, func() {
				s.mu.Lock()
				s.executeNow()
				s.mu.Unlock()
			})
		} else {
			// Execute immediately after this call completes
			go func() {
				s.mu.Lock()
				s.executeNow()
				s.mu.Unlock()
			}()
		}
	}

	s.mu.Unlock()

	// Wait for result
	select {
	case result := <-resultCh:
		return result.value, result.err
	case <-ctx.Done():
		var zero R
		return zero, ctx.Err()
	}
}

// executeNow executes the current batch immediately.
// Must be called with lock held.
func (s *BatchScheduler[P, R]) executeNow() {
	if len(s.pending) == 0 {
		return
	}

	// Cancel timer if running
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}

	// Capture current batch
	batch := s.pending
	s.pending = make([]pendingItem[P, R], 0)

	// Extract args
	args := make([]P, len(batch))
	for i, item := range batch {
		args[i] = item.arg
	}

	// Execute in background
	go func() {
		results, err := s.config.Fn(context.Background(), args)

		// Send results to all waiters
		for i, item := range batch {
			var result batchResult[R]
			if err != nil {
				result.err = err
			} else if i < len(results) {
				result.value = results[i]
			}
			item.result <- result
		}
	}()
}

// Flush executes any pending batch immediately and clears the scheduler.
func (s *BatchScheduler[P, R]) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executeNow()
}

// Clear removes this scheduler from the cache.
func (s *BatchScheduler[P, R]) Clear() {
	schedulerCache.Delete(s.config.ID)
}

// ClearSchedulerCache clears all cached schedulers.
func ClearSchedulerCache() {
	schedulerCache.Range(func(key, value interface{}) bool {
		schedulerCache.Delete(key)
		return true
	})
}
