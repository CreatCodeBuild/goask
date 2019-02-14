package log

import (
	"log"
	"os"
)

type Logger struct{}

func (l *Logger) Error(err error) {
	log.Printf("%+v\n", err)
}

func (l *Logger) ErrorExit(err error) {
	log.Printf("%+v\n", err)
	os.Exit(1)
}
