package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type Logger struct {
	Level     int // 0: debug, 1: info, 2: error 3: warnning
	TimeStamp string
	Message   string
}
type LogFile struct {
	file     *os.File
	fileName string
	mu       sync.RWMutex
	enabled  bool
}

func NewLogFile(fileName string) (*LogFile, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &LogFile{
		file:     file,
		fileName: fileName,
		enabled:  true,
	}, nil
}
func (lf *LogFile) WriteLog(log Logger) error {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	if !lf.enabled {
		return errors.New("log file is disabled")
	}
	level := map[int]string{
		0: "DEBUG",
		1: "INFO",
		2: "WARNING",
		3: "ERROR",
		4: "FATAL",
	}
	logEntry := fmt.Sprintf("%s [%s] %s\n", log.TimeStamp, level[log.Level], log.Message)
	_, err := lf.file.WriteString(logEntry)
	return err
}
func (lf *LogFile) ReadLog() error {
	lf.mu.RLock()
	defer lf.mu.RUnlock()
	if !lf.enabled {
		return errors.New("log file is disabled")
	}
	content, err := os.ReadFile(lf.fileName)
	if err != nil {
		return err
	}
	fmt.Println(string(content))
	return nil
}
func (lf *LogFile) BanchedWriteLog(log []Logger) error {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	if !lf.enabled {
		return errors.New("log file is disabled")
	}
	var logEntries string
	for _, l := range log {
		logEntries += fmt.Sprintf("%s [%d] %s\n", l.TimeStamp, l.Level, l.Message)
	}
	_, err := lf.file.WriteString(logEntries)
	return err
}
func (lf *LogFile) Close() error {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	if !lf.enabled {
		return errors.New("log file is disabled")
	}
	return lf.file.Close()
}
