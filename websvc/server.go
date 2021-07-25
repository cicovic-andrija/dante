package websvc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/cicovic-andrija/dante/db"
	"github.com/cicovic-andrija/dante/util"
)

const (
	HTTPClientTimeout = 15 * time.Second
)

type server struct {
	// server
	name      string
	shutdownC chan struct{}

	// logging
	log *logstruct

	// http server
	httpServer *http.Server
	httpWg     *sync.WaitGroup
	router     *mux.Router

	// clients
	httpClient *http.Client
	dbcli      *db.Client

	// timer tasks
	timerTasks  []*timerTask
	timerTaskWg *sync.WaitGroup
}

func (s *server) init() error {
	var err error

	if s.name, err = util.RandHexString(6); err != nil {
		return err
	}

	s.shutdownC = make(chan struct{})

	s.log = &logstruct{}
	if err = s.log.init(s.name); err != nil {
		return err
	}

	s.log.info("server %s (version: %s)", s.name, version)
	s.log.info("configuration: %s", cfg.Path())
	s.log.info("environment: %s", cfg.Env)

	s.httpInit()
	s.httpClientInit()
	s.dbinit()

	s.timerTaskWg = &sync.WaitGroup{}
	s.timerTasks = []*timerTask{
		{name: "get-credits", execute: s.getCredits, period: 10 * time.Second, log: s.log},
	}

	return nil
}

func (s *server) run() {
	// HTTP
	go s.httpRun()

	// Timer tasks
	s.log.info("[main] starting timer tasks ...")
	for _, task := range s.timerTasks {
		task.run(s.timerTaskWg)
		s.timerTaskWg.Add(1)
	}

	s.log.info("[main] boot sequence completed")
}

func (s *server) signalShutdown() {
	close(s.shutdownC)
}

func (s *server) shutdown() {
	s.log.info("[main] shutting down the server ...")

	for _, task := range s.timerTasks {
		task.stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx) // FIXME: Handle errors returned by Shutdown.

	s.timerTaskWg.Wait()
	s.httpWg.Wait()

	s.log.info("[main] server stopped.")
	s.log.finalize()
}

func (s *server) dbinit() error {
	// database client
	s.dbcli = db.NewClient(&cfg.Influx)

	if created, err := db.EnsureOrganization(s.dbcli); err != nil {
		return fmt.Errorf("failed to ensure influxdb org %q is created: %v", s.dbcli.Organization, err)
	} else if created {
		s.log.info("[main] influxdb org %q created", s.dbcli.Organization)
	}
	return nil
}

func (s *server) interruptHandler() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	s.log.info("[interrupt] signaling shutdown ...")
	s.signalShutdown()
}
