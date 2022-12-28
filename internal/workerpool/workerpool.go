package workerpool

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"
)

type Job interface {
	Run(ctx context.Context) error
}

func RunPool(ctx context.Context, size int, jobs chan Job) error {
	gr, ctx := errgroup.WithContext(ctx)
	for i := 0; i < size; i++ {
		gr.Go(func() error {
			for {
				select {
				case job := <-jobs:
					err := job.Run(ctx)
					if err != nil {
						fmt.Printf("Job error: %s \n", err)
						return err
					}
				case <-ctx.Done():
					fmt.Println("Context canceled")
					return ctx.Err()
				}
			}
		})
	}

	return gr.Wait()
}

func RunPoolV2(ctx context.Context, size int, jobs chan Job) error {
	token := make(chan struct{}, size)

	for {
		select {
		case job := <-jobs:
			log.Printf("got job to handle")

			go func() {
				token <- struct{}{}
				defer func() { <-token }()
				err := job.Run(ctx)
				if err != nil {
					fmt.Printf("Job error: %s \n", err)
					return
				}
				log.Printf("job handled")
			}()
		case <-ctx.Done():
			fmt.Println("Context canceled")
			return ctx.Err()
		}
	}
}
