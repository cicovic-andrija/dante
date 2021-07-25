package websvc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type logger struct {
	file    *os.File
	backend *log.Logger
}

type logstruct struct {
	infoLogger  *logger
	errorLogger *logger
}

func (L *logstruct) info(format string, v ...interface{}) {
	L.infoLogger.backend.Printf(format, v...)
}

func (L *logstruct) err(format string, v ...interface{}) {
	L.errorLogger.backend.Printf(format, v...)
}

func (L *logstruct) fatal(format string, v ...interface{}) {
	L.errorLogger.backend.Fatalf(format, v...)
}

func (L *logstruct) init(serverId string) error {
	var (
		err error
	)

	initOne := func(logType string, logPrefix string, symlink string) (*logger, error) {
		file, path, err := createLogFile(serverId, logType)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %v", err)
		}

		backend := log.New(
			file,
			logPrefix,
			log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lmsgprefix,
		)

		// try to create or update a symlink to the log file
		// fail silently
		updateSymlink(filepath.Join(cfg.Log.Dir, symlink), path)

		return &logger{file: file, backend: backend}, nil
	}

	if L.infoLogger, err = initOne("inf", "I: ", "info"); err != nil {
		return err
	}

	if L.errorLogger, err = initOne("err", "E: ", "error"); err != nil {
		return err
	}

	return nil
}

// FIXME: Handle errors.
func (L *logstruct) finalize() {
	if L.infoLogger.file != nil {
		L.infoLogger.file.Close()
	}
	if L.errorLogger.file != nil {
		L.errorLogger.file.Close()
	}
}

func createLogFile(prefix string, suffix string) (*os.File, string, error) {
	// generate file name
	filename := fmt.Sprintf(
		"%s_%s_%s.log",
		prefix,
		time.Now().UTC().Format("D20060102_T150405"),
		suffix,
	)

	path := filepath.Join(cfg.Log.Dir, filename)
	file, err := os.Create(path)
	if err != nil {
		return nil, "", err
	}

	return file, path, nil
}

// FIXME: Handle errors.
func updateSymlink(symlink string, target string) {
	if _, err := os.Lstat(symlink); err == nil {
		if err := os.Remove(symlink); err == nil {
			os.Symlink(target, symlink)
		}
	} else if os.IsNotExist(err) {
		os.Symlink(target, symlink)
	}
}
