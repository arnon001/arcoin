package main

import (
	"./arcoin/cli"
	"./arcoin/core"
	"./arcoin/p2p"
	"flag"
)

func main() {
	nodeID := flag.String("node", "3000", "Node ID")
	flag.Parse()

	bc := core.NewBlockchain()
	node := p2p.StartNode(*nodeID, "")
	
	cli := cli.CLI{BC: bc, Node: node}
	cli.Run()
}