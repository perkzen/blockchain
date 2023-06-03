package cli

import (
	"blockchain/pkg/network"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

type CommandLine struct {
}

func (cli *CommandLine) listCommands() {
	fmt.Println("Usage: ")
	//fmt.Println("create wallet")
	fmt.Println("start -port --> start blockchain server")
	//fmt.Println("print chain")
	//fmt.Println("send")
	//fmt.Println("get balance")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.listCommands()
		runtime.Goexit()
	}

}

func (cli *CommandLine) StartServer(port uint16) {
	app := network.NewBlockchainServer(port)
	app.Run()
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	startNodeCmd := flag.NewFlagSet("start", flag.ExitOnError)
	port := startNodeCmd.Uint("port", 3000, "TCP Port Number for Blockchain server")

	switch os.Args[1] {
	case "start":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		cli.listCommands()
		runtime.Goexit()
	}

	if startNodeCmd.Parsed() {
		cli.StartServer(uint16(*port))
	}
}
