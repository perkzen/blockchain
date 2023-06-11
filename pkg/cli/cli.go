package cli

import (
	"blockchain/pkg/network"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

type CommandLine struct {
}

const (
	START_NODE = "start-node"
	SEND       = "send"
)

func (cli *CommandLine) listCommands() {
	fmt.Println("Usage: ")
	fmt.Printf("%s -port --> start blockchain server\n", START_NODE)
	fmt.Printf("%s -to -amount --> send amount from address to address\n", SEND)

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

func (cli *CommandLine) send(to string, amount float64) {
	body := []byte(fmt.Sprintf(`{"recipient": "%s", "amount": %f}`, to, amount))
	res, err := http.Post("http://localhost:3000/transaction", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Panic(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Panic("ERROR: Failed to send transaction")
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(bodyBytes))
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	startNodeCmd := flag.NewFlagSet(START_NODE, flag.ExitOnError)
	port := startNodeCmd.Uint("port", 3000, "TCP Port Number for Blockchain server")

	sendCmd := flag.NewFlagSet(SEND, flag.ExitOnError)
	to := sendCmd.String("to", "", "Recipient of the transaction")
	amount := sendCmd.Float64("amount", 0, "Amount to send")

	switch os.Args[1] {
	case START_NODE:
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case SEND:
		err := sendCmd.Parse(os.Args[2:])
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

	if sendCmd.Parsed() {
		cli.send(*to, *amount)
	}
}
