package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/yankeguo/fastci"
	"github.com/yankeguo/rg"
)

const (
	fileStdin = "-"
)

func main() {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}()
	defer rg.Guard(&err)

	var (
		optFile string
	)
	flag.StringVar(&optFile, "f", fileStdin, "fastci script file to read, - for stdin")
	flag.Parse()

	var f *os.File
	if optFile == fileStdin {
		f = os.Stdin
	} else {
		f = rg.Must(os.OpenFile(optFile, os.O_RDONLY, 0))
		defer f.Close()
	}

	rg.Must0(fastci.NewRunner().Execute(context.Background(), f))
}
