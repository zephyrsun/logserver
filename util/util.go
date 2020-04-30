package util

import (
	"log"
)

func Fatal(s interface{}) {
	if s == nil {
		return
	}

	log.Fatal(s)
}

func Error(format string, err error) {
	if err != nil {
		log.Print(format, err)
	}
}

func Print(format string, v ...interface{}) {
	log.Printf(format, v...)
}
