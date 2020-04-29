package util

import (
	"log"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Error(format string, err error) {
	if err != nil {
		log.Print(format, err)
	}
}

func Print(format string, v ...interface{}) {
	log.Printf(format, v...)
}
