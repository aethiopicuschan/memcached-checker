package main

import (
	"os"

	"github.com/aethiopicuschan/memcached-checker/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
