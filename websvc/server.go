package websvc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cicovic-andrija/dante/ripe"
	"github.com/cicovic-andrija/dante/util"
)

type server struct {
	name        string
	httpClient  *http.Client
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
	s.inflog.Printf("name: %s, version: %s, build: %s, env: %s",
		s.name, "0.1", "0.0", cfg.Env)

	s.inflog.Printf("configuration: %s", cfg.Path())

	s.httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	return nil
}

func (s *server) periodicCreditCheck() {
	fmt.Println("Checking credit..")
	// FIXME: handle error
	req, _ := http.NewRequest(http.MethodGet, "https://atlas.ripe.net:443/api/v2/credits/", nil)
	req.Header.Set("Authorization", "Key "+cfg.Auth.Key)
	res, err := s.httpClient.Do(req)
	if err != nil {
		fmt.Println("error while requesting")
		return
	}
	c := ripe.Credit{}
	err = json.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		fmt.Println("error while decoding")
		return
	}
	fmt.Println(c.CurrentBalance)
}

func (s *server) run() {

}

func (s *server) shutdown() {

}

// FIXME: Handle errors and file closure
func (s *server) initLogs() {
	s.logRotation = 1

	path := filepath.Join(cfg.Log.Dir, util.GenerateLogName(s.name, s.logRotation, "info"))
	file, _ := os.Create(path)
	s.inflog = log.New(file, "I: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)

	path = filepath.Join(cfg.Log.Dir, util.GenerateLogName(s.name, s.logRotation, "error"))
	file, _ = os.Create(path)
	s.errlog = log.New(file, "E: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
}
