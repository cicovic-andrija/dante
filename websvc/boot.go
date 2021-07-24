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

func Start() {
	var err error
	if cfg, err = conf.Load(); err != nil {
		die(err)
	}

	srv := &server{}
	err = srv.init()
	if err != nil {
		die(fmt.Errorf("server initialization failed: %v", err))
	}

	srv.run()
}
