package main

import (
	"os"

	"github.com/antooni/go-blockchain/cli"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
