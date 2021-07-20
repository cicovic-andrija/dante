package main

import (
	"flag"

	"github.com/cicovic-andrija/dante/websvc"
)

func main() {
	flag.Parse()
	websvc.Start()
}
