package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("args is empty")
	}

	environment, err := ReadDir(args[0])
	if err != nil {
		log.Fatalf("failed to ReadDir: %+v", err)
	}

	returnCode := RunCmd(args[1:], environment)

	os.Exit(returnCode)
}
