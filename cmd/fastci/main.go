package main

import (
	"context"
	"log"
	"os"

	"github.com/yankeguo/fastci"
)

func main() {
	if err := fastci.NewRunner().Execute(context.Background(), os.Stdin); err != nil {
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}
}
