package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cicovic-andrija/dante/config"
)

func main() {
	var (
		configPath = flag.String("config", "", "Config file path")
		err        error
	)

	flag.Parse()

	if _, err = config.Load(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
