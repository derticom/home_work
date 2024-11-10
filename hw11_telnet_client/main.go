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
	defer func(client TelnetClient) {
		err := client.Close()
		if err != nil {
			log.Println("failed to client.Close:", err)
		}
	}(client)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	errChan := make(chan error, 1)

	go func() {
		err := client.Send()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				errChan <- errors.WithMessage(err, "failed to client.Send")
				return
			}

			return
		}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("connection closed by peer")
			} else {
				errChan <- errors.WithMessage(err, "failed to client.Receive")
			}

			return
		}
	}()

	select {
	case <-sigChan:
		return
	case err := <-errChan:
		log.Printf("critical error, %+v", err)
	}
}
