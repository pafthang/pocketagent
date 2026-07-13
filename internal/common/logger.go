package common

import (
	"log"
	"os"
)

// Logger provides structured logging

func NewLogger(service string) *log.Logger {
	return log.New(os.Stdout, "["+service+"] ", log.LstdFlags|log.Lshortfile)
}
