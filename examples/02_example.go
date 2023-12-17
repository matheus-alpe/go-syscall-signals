package examples

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const workerPoolSize = 4

type Consumer struct {
	ingestChan <-chan int
	jobsChan   chan int
}

func (c *Consumer) Run(ctx context.Context) {
	fmt.Println("Consumer: start")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Consumer: stop")
			close(c.jobsChan)
			return
		case value := <-c.ingestChan:
			// Process the received value (for example, send it as a job).
			fmt.Printf("Consumer: Received value %d. Sending as a job.\n", value)
			c.jobsChan <- value // Sending the value as a job (processing example).
		}
	}
}

func (c *Consumer) Worker(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for value := range c.jobsChan {
		time.Sleep(time.Second)
		fmt.Println("Worker:", id, "with job:", value)
	}
}

type Producer struct {
	ingestChan chan<- int
}

func (p *Producer) Run(ctx context.Context) {
	fmt.Println("Producer: start")

	for i := 1; i < 100; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Producer: stop")
			return
		default:
			p.ingestChan <- i
		}
	}
}

func Example02() {
	ingestChan := make(chan int, 1)

	// Simulate some input
	producer := Producer{ingestChan: ingestChan}

	// Simulate some consumer
	consumer := Consumer{
		ingestChan: ingestChan,
		jobsChan:   make(chan int, workerPoolSize),
	}

	// Set up cancellation context and wait group
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Start control loop with cancel context
	go producer.Run(ctx)
	go consumer.Run(ctx)

	// Start workers
	wg.Add(workerPoolSize)
	for i := 0; i < workerPoolSize; i++ {
		go consumer.Worker(wg, i)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // blocks here until interrupted

	fmt.Println("=== Shutdown received ===")
	cancelFunc() // Signal cancel to context
	wg.Wait()    // wait for workers

	fmt.Println("All workers done, shutting down!")
}
