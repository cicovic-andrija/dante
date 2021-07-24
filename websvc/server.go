package websvc

import (
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/cicovic-andrija/dante/util"
)

type server struct {
	name           string
	httpClient     *http.Client
	influxDbClient influxdb2.Client
	timerTasks     []*timerTask
	log            *logstruct
}

func (s *server) init() error {
	var err error

	if s.name, err = util.RandHexString(6); err != nil {
		return err
	}

	s.log = &logstruct{}
	if err = s.log.init(s.name); err != nil {
		return err
	}

	s.httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	s.timerTasks = []*timerTask{
		{name: "get-credits", task: s.getCredits, period: 5 * time.Second, log: s.log},
	}

	// finalize init.
	s.log.info("server %s (version: %s)", s.name, Version)
	s.log.info("configuration: %s", cfg.Path())
	s.log.info("environment: %s", cfg.Env)
	s.log.info("boot sequence completed")
	return nil
}

func (s *server) run() {
	s.log.info("(main) starting timer tasks ...")
	for _, task := range s.timerTasks {
		task.run()
	}

	time.Sleep(60 * time.Second)
	s.shutdown()
}

func (s *server) shutdown() {
	s.log.info("(main) stopping timer tasks ...")
	for _, task := range s.timerTasks {
		task.stop()
	}

	s.log.finalize()
}
