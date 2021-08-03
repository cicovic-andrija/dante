package websvc

import (
	"fmt"
	"os"

	"github.com/cicovic-andrija/dante/conf"
)

var cfg *conf.Config

func die(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

// Start initializes and starts the web service.
func Start() {
	var err error

	if cfg, err = conf.Load(); err != nil {
		die(err)
	}

	// initialize a server
	srv := &server{}
	if err = srv.init(); err != nil {
		die(fmt.Errorf("server initialization failed: %v", err))
	}

	// run the server itself
	srv.run()

	// run an interrupt handler in a separate thread
	go srv.interruptHandler()

	// wait for the shutdown request and shut down the server
	<-srv.shutdownC
	srv.shutdown()
}
