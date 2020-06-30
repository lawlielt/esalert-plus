package log

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type Writer struct {
	LogDir    string
	LogLevels []log.Level
}

func NewWriter(logDir string, logLevels []log.Level) (*Writer, error) {
	fstat, err := os.Stat(logDir)
	if err != nil || !fstat.IsDir(){
		return nil, fmt.Errorf("%s, dir not exists", logDir)
	}
	writer := new(Writer)
	writer.LogDir = logDir
	writer.LogLevels = logLevels
	return writer, nil
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (writer *Writer) Fire(entry *log.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}
	fileName := filepath.Join(writer.LogDir, time.Now().Format("20060102") + ".log")
	file, err := os.OpenFile(fileName, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(line)
	return err
}

// Levels define on which log levels this hook would trigger
func (writer *Writer) Levels() []log.Level {
	return writer.LogLevels
}
