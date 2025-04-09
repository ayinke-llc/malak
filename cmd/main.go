package main

import (
	"log"
	"os"

	"github.com/ayinke-llc/malak/cmd/cli"
)

func main() {
	os.Setenv("TZ", "")

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
