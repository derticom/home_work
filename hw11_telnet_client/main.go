package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors" //nolint: depguard // import is necessary
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 2 {
		log.Fatal("host and port required")
	}

	host := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(host, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal("failed to client.Connect:", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Println("failed to client.Close:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	errChan := make(chan error, 1)

	go func() {
		if err := client.Send(); err != nil && !errors.Is(err, io.EOF) {
			errChan <- errors.WithMessage(err, "failed to client.Send")
		} else {
			sigChan <- syscall.SIGINT
		}
	}()

	go func() {
		if err := client.Receive(); err != nil && !errors.Is(err, io.EOF) {
			errChan <- errors.WithMessage(err, "failed to client.Receive")
		} else {
			sigChan <- syscall.SIGINT
		}
	}()

	select {
	case <-sigChan:
		log.Println("shut down")
	case err := <-errChan:
		log.Printf("critical error, %+v", err)
	}
}
