package log

import (
	"log"
)

type Logger struct{}

func (l *Logger) Error(err error) {
	log.Printf("%+v\n", err)
}
