package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Simulate reading, validating, and writing a file
func processFile(ctx context.Context, workerID int, file string) bool {
	// Check for cancellation before starting
	select {
	case <-ctx.Done():
		fmt.Printf("[Worker %d] Skipping %s due to cancellation.\n", workerID, file)
		return false
	default:
	}

	fmt.Printf("[Worker %d] Reading %s...\n", workerID, file)
	time.Sleep(100 * time.Millisecond) // Simulate read

	fmt.Printf("[Worker %d] Validating %s...\n", workerID, file)
	time.Sleep(50 * time.Millisecond) // Simulate validation

	fmt.Printf("[Worker %d] Writing %s to DB...\n", workerID, file)
	time.Sleep(150 * time.Millisecond) // Simulate DB write

	fmt.Printf("[Worker %d] Done processing %s.\n", workerID, file)
	return true
}

func main() {
	// Create a context that can be cancelled (for graceful shutdown)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)

	// Notify sigChan if SIGINT (Ctrl+C) or SIGTERM is received
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a Goroutine to handle signals
	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal: %s. Shutting down gracefully...\n", sig)
		cancel() // Trigger context cancellation
	}()

	// Simulate 100 file names
	files := make([]string, 100)
	for i := 0; i < 100; i++ {
		files[i] = fmt.Sprintf("file_%03d.txt", i+1)
	}

	numWorkers := 10                        // Number of worker Goroutines
	jobs := make(chan string, len(files))   // Buffered channel for file jobs
	var wg sync.WaitGroup                   // WaitGroup to wait for workers to finish

	var processedCount int
	var mu sync.Mutex                       // Mutex to protect processedCount

	// Start worker pool
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					// Exit worker if cancelled
					fmt.Printf("[Worker %d] Exiting due to cancellation.\n", workerID)
					return
				case file, ok := <-jobs:
					if !ok {
						// Channel closed, no more jobs
						return
					}
					if processFile(ctx, workerID, file) {
						// Safely increment processedCount
						mu.Lock()
						processedCount++
						mu.Unlock()
					}
				}
			}
		}(w)
	}

	// Feed jobs to the workers
jobFeed:
	for _, file := range files {
		select {
		case <-ctx.Done():
			fmt.Println("Job feeding stopped due to cancellation.")
			break jobFeed
		default:
			jobs <- file
		}
	}
	close(jobs) // No more jobs to send

	// Wait for all workers to finish
	wg.Wait()

	fmt.Printf("All done. Total files processed: %d\n", processedCount)
}
