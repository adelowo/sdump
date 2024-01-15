package main

import (
	"log"

	"github.com/adelowo/sdump/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
