package main

import (
	"fmt"
	"os"

	"github.com/cicovic-andrija/rac/config"
)

var (
	g_cfg *config.Config
)

func main() {
	var (
		err error
	)

	if g_cfg, err = config.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
