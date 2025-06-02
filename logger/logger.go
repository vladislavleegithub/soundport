package logger

import (
	"log"
	"os"
	"sync"
)

type notFoundLogger struct {
	filename string
	*log.Logger
}

type globalLogger struct {
	filename string
	*log.Logger
}

var (
	glbLogger *globalLogger
	nfLogger  *notFoundLogger
	glbOnce   sync.Once
	nfOnce    sync.Once
)

// start loggeando
func GetInstance() *globalLogger {
	glbOnce.Do(func() {
		glbLogger = createLogger("/tmp/soundport.log")
	})
	return glbLogger
}

func GetNotFoundLogInstance() *notFoundLogger {
	nfOnce.Do(func() {
		nfLogger = createNfLogger("/tmp/sp_notfound.log")
	})
	return nfLogger
}

func createNfLogger(fname string) *notFoundLogger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	return &notFoundLogger{
		filename: fname,
		Logger:   log.New(file, "Log: ", log.LstdFlags),
	}
}

func createLogger(fname string) *globalLogger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	return &globalLogger{
		filename: fname,
		Logger:   log.New(file, "Error: ", log.Lshortfile|log.LstdFlags),
	}
}
