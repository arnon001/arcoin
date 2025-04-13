package cli

import (
	"arcoin/core"
	"arcoin/p2p"
	"arcoin/wallet"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	BC   *core.Blockchain
	Node *p2p.Node
}

func (cli *CLI) Run() {
	cli.printUsage()
	
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		switch cmd {
		case "createwallet":
			cli.createWallet()
		case "getbalance":
			cli.getBalance()
		case "send":
			cli.send()
		case "mine":
			cli.mine()
		case "startnode":
			cli.startNode()
		case "exit":
			os.Exit(0)
		default:
			cli.printUsage()
		}
	}
}

func (cli *CLI) createWallet() {
	w := wallet.NewWallet()
	fmt.Printf("New address: %s\n", w.Address())
}

func (cli *CLI) getBalance() {
	// Implement balance calculation from UTXO set
}

func (cli *CLI) send() {
	// Implement transaction creation and signing
}

func (cli *CLI) mine() {
	// Implement mining command
}

func (cli *CLI) startNode() {
	// Start P2P node
}

func (cli *CLI) printUsage() {
	fmt.Println("Commands:")
	fmt.Println(" createwallet - Generate new wallet")
	fmt.Println(" getbalance - Check balance")
	fmt.Println(" send - Send coins")
	fmt.Println(" mine - Mine new block")
	fmt.Println(" startnode - Start P2P node")
	fmt.Println(" exit - Exit program")
}