package logger

import (
	"log"
	"os"
	"sync"
)

type globalLogger struct {
	filename string
	*log.Logger
}

var (
	logger *globalLogger
	once   sync.Once
)

// start loggeando
func GetInstance() *globalLogger {
	once.Do(func() {
		logger = createLogger("/tmp/soundport.log")
	})
	return logger
}

func createLogger(fname string) *globalLogger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	return &globalLogger{
		filename: fname,
		Logger:   log.New(file, "Error: ", log.Lshortfile|log.LstdFlags),
	}
}
