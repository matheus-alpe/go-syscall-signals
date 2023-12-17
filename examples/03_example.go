package examples

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Settings struct {
	config string
	sigusr chan os.Signal
	done chan struct{}
	sync.Mutex
}

func NewSettings() *Settings {
	s := &Settings{sigusr: make(chan os.Signal, 1)}
	s.readConfigFromFile()
	go s.HotReloadConfig()
	return s
}

func (s *Settings) Config() string {
	s.Lock()
	cnf := s.config
	s.Unlock()
	return cnf
}

func (s *Settings) readConfigFromFile() {
	content, err := os.ReadFile("settings.conf")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s.Lock()
	s.config = string(content)
	fmt.Println(s.config)
	s.Unlock()
}

func (s *Settings) HotReloadConfig() {
	for {
		<-s.sigusr // this line will block until the signal is received
		s.readConfigFromFile()
	}
}

func Example03() {
	fmt.Printf("Process pid: %v\n", os.Getpid())

	settings := NewSettings()
	// to simulate, execute on terminal while the program is running
	// $ kill -SIGUSR1 <process-pid>
	signal.Notify(settings.sigusr, syscall.SIGUSR1)
	<-settings.done
}
