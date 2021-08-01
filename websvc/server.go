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

const (
	serverNameHexLength = 6
	taskPoolInitCap     = 32
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

	// measurements
	measCache *measurementCache

	// timer tasks
	taskManager *timerTaskManager
}

func (s *server) init() error {
	var err error

	if s.name, err = util.RandHexString(serverNameHexLength); err != nil {
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

	s.measCache = newMeasurementCache()

	s.taskManager = &timerTaskManager{
		tasks: make([]*timerTask, 0, taskPoolInitCap),
		wg:    &sync.WaitGroup{},
	}

	// default, always-running timer tasks
	s.taskManager.addTask("get-credits", s.getCredits, 5*time.Minute, s.log)
	s.taskManager.addTask("probe-database", s.probeDatabase, 10*time.Minute, s.log)

	return nil
}

func (s *server) run() {
	// HTTP
	go s.runHTTP()

	// Timer tasks
	s.log.info("[main] starting timer tasks ...")
	s.taskManager.runAll()
}

func (s *server) signalShutdown() {
	close(s.shutdownC)
}

func (s *server) shutdown() {
	s.log.info("[main] shutting down the server ...")

	s.taskManager.stopAll()
	s.database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx) // TODO: Handle error returned by Shutdown.

	s.taskManager.wg.Wait()
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

	if org, err := s.database.EnsureOrganization(cfg.Influx.Organization); err != nil {
		return formatError(err)
	} else {
		s.database.Org = org
	}

	if bck, err := s.database.EnsureBucket(db.OperationalDataBucket); err != nil {
		return formatError(err)
	} else {
		s.database.OperDataBucket = bck
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
