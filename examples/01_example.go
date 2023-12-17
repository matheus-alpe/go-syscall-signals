package examples

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Example01() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	fmt.Printf("Waiting for signal, pid: %v\n", os.Getpid())

	sig := <-sigChan
	fmt.Println("Program", sig)
	fmt.Println("Program finished")
}
