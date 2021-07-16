package websvc

import (
	"fmt"
	"os"

	"github.com/cicovic-andrija/dante/conf"
)

var cfg *conf.Config

func die(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}

func Start(confPath string) {
	var err error
	if cfg, err = conf.Load(confPath); err != nil {
		die(err)
	}

	srv := &server{}
	srv.init()
	srv.run()
}
