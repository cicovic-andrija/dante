package main

import (
	"flag"

	"github.com/cicovic-andrija/dante/websvc"
)

func main() {
	conf := flag.String("conf", "", "Config file path")
	flag.Parse()
	websvc.Start(*conf)
}
