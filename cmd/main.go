package main

import (
	"blockchain/pkg/cli"
	"log"
	"os"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}
