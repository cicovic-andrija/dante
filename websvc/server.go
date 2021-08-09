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

// server specifies the values and structures needed
// during the service's run time.
type server struct {
	// server/control
	name      string
	shutdownC chan struct{}

	// logging
	log *logstruct

	// http server
	httpServer *http.Server
	httpWg     *sync.WaitGroup
	router     http.Handler

	// http client
	httpClient *http.Client

	// database objects
	database *db.Client
	mmd      []db.MeasurementMetadata

	// measurements
	measCache *measurementCache
	probeInfo *probeTable

	// timer tasks
	taskManager   *timerTaskManager
	taskManagerWg *sync.WaitGroup
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

	s.log.info("[main] server %s (version: %s)", s.name, version)
	s.log.info("[main] configuration: %s", cfg.Path())
	s.log.info("[main] environment: %s", cfg.Env)

	s.httpInit()

	if err = s.dbinit(); err != nil {
		return err
	}

	s.measCache = newMeasurementCache()
	s.probeInfo = newProbeTable()

	s.taskManager = newTimerTaskManager()

	// default, always-running timer tasks
	s.taskManager.addTask("get-credits", s.getCredits, 5*time.Minute, s.log)
	s.taskManager.addTask("probe-database", s.probeDatabase, 10*time.Minute, s.log)

	return nil
}

func (s *server) run() {
	// HTTP server
	go s.runHTTP()

	// Timer tasks
	s.log.info("[main] starting timer tasks ...")
	s.taskManagerWg = &sync.WaitGroup{}
	s.taskManagerWg.Add(1)
	s.taskManager.run(s.taskManagerWg) // run task manager itself
	s.taskManager.runAllTasks()

	// Restore thread
	go s.restore()
}

func (s *server) signalShutdown() {
	close(s.shutdownC)
}

func (s *server) shutdown() {
	s.log.info("[main] shutting down the server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx) // TODO: Handle error returned by Shutdown.
	s.httpWg.Wait()

	s.taskManager.stopAll()
	s.taskManager.stop()

	s.database.Close()

	s.taskManager.wg.Wait()
	s.taskManagerWg.Wait()

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

	if bck, err := s.database.EnsureBucket(db.SystemBucket); err != nil {
		return formatError(err)
	} else {
		s.database.SystemBucket = bck
	}

	if mmd, err := s.database.QueryMeasurementMetadata(); err != nil {
		return formatError(err)
	} else {
		s.mmd = mmd
	}

	return nil
}

// Probably no one will be shutting down the server with a signal,
// but add a signal handler anyway, it's just 9 lines of code.
func (s *server) interruptHandler() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	s.log.info("[interrupt] signaling shutdown ...")
	s.signalShutdown()
}
