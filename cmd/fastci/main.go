package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/yankeguo/fastci"
	"github.com/yankeguo/rg"
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

	p := fastci.NewPipeline()
	script := rg.Must(io.ReadAll(os.Stdin))
	err = p.Do(context.Background(), string(script))
}
