package logger

import (
	"io"
	"log"
	"os"
	"runtime/debug"
)

var handle HandleFunc

type Option struct {
	writer io.Writer
	Handle HandleFunc
}

type HandleFunc func(string)

func Init(option Option) {
	handle = option.Handle

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(io.MultiWriter(option.writer, os.Stdout))
}

func Error(em string) {
	log.Println(em)
	log.Println(string(debug.Stack()))
	if handle != nil {
		handle(em)
	}
}
