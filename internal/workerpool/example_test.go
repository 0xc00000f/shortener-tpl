package workerpool

import (
	"context"
	"log"
)

// Example worker run
func ExampleRunPool() {
	concurrency := 10
	jobsCh := make(chan Job, concurrency)

	err := RunPool(context.Background(), concurrency, jobsCh)
	if err != nil {
		log.Printf("runpool err: %v", err)
	}
}
