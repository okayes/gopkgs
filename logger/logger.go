package logger

import (
	"io"
	"log"
	"os"
	"runtime/debug"
)

var handle HandleFunc

type Option struct {
	FilePath, FileName string
	Handle             HandleFunc
}

type HandleFunc func(string)

func Init(option Option) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	handle = option.Handle

	if _, err := os.Stat(option.FilePath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(option.FilePath, os.ModePerm); err != nil {
				log.Println(err)
				return
			}
		} else {
			log.Println(err)
			return
		}
	}

	file, err := os.Create(option.FilePath + option.FileName)
	if err != nil {
		log.Println(err)
		return
	}

	log.SetOutput(io.MultiWriter(file, os.Stdout))
}

func Error(em string) {
	log.Println(em)
	log.Println(string(debug.Stack()))
	if handle != nil {
		handle(em)
	}
}
