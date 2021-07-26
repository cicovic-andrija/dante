package websvc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/cicovic-andrija/dante/db"
	"github.com/cicovic-andrija/dante/util"
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
	router     http.Handler

	// clients
	httpClient *http.Client
	database   *db.Client

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

	if err = s.dbinit(); err != nil {
		return err
	}

	s.timerTaskWg = &sync.WaitGroup{}
	s.timerTasks = []*timerTask{
		{name: "get-credits", execute: s.getCredits, period: 10 * time.Minute, log: s.log},
		{name: "probe-database", execute: s.probeDatabase, period: 30 * time.Minute, log: s.log},
	}

	return nil
}

func (s *server) run() {
	// HTTP
	go s.runHTTP()

	// Timer tasks
	s.log.info("[main] starting timer tasks ...")
	for _, task := range s.timerTasks {
		task.run(s.timerTaskWg)
		s.timerTaskWg.Add(1)
	}
}

func (s *server) signalShutdown() {
	close(s.shutdownC)
}

func (s *server) shutdown() {
	s.log.info("[main] shutting down the server ...")

	for _, task := range s.timerTasks {
		task.stop()
	}

	s.database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx) // FIXME: Handle error returned by Shutdown.

	s.timerTaskWg.Wait()
	s.httpWg.Wait()

	s.log.info("[main] server stopped.")
	s.log.finalize()
}

func (s *server) dbinit() error {
	// database client
	s.database = db.NewClient(&cfg.Influx)

	formatError := func(err error) error {
		return fmt.Errorf("failed to init database: %v", err)
	}

	if err := s.database.EnsureOrganization(cfg.Influx.Organization); err != nil {
		return formatError(err)
	}

	if err := s.database.EnsureBuckets(); err != nil {
		return formatError(err)
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
