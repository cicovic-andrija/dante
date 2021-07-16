package websvc

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cicovic-andrija/dante/util"
)

type server struct {
	name        string
	inflog      *log.Logger
	errlog      *log.Logger
	logRotation int
}

func (s *server) init() error {
	var err error

	if s.name, err = util.GenerateHexString(6); err != nil {
		return err
	}

	s.initLogs()

	// finalize init
	s.inflog.Printf("name: %s, version: %s, build: %s, env: %s", s.name, "0.1", "0.0", cfg.Env)
	s.inflog.Printf("configuration: %s", cfg.Path())
	return nil
}

func (s *server) run() {

}

func (s *server) shutdown() {

}

func (s *server) initLogs() {
	// FIXME: Handle errors and file closure
	s.logRotation = 1

	path := filepath.Join(cfg.Log.Dir, util.GenerateLogName(s.name, s.logRotation, "info"))
	file, _ := os.Create(path)
	s.inflog = log.New(file, "I: " /* prefix */, log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	if cfg.Log.InfoVar != "" {
		os.Setenv(cfg.Log.InfoVar, path)
	}

	path = filepath.Join(cfg.Log.Dir, util.GenerateLogName(s.name, s.logRotation, "error"))
	file, _ = os.Create(path)
	s.errlog = log.New(file, "E: " /* prefix */, log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	if cfg.Log.ErrorVar != "" {
		os.Setenv(cfg.Log.ErrorVar, path)
	}
}
