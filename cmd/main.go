package main

import (
	"blockchain/pkg/network"
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 3000, "TCP Port Number for Blockchain server")
	flag.Parse()
	app := network.NewBlockchainServer(uint16(*port))
	app.Run()
}
