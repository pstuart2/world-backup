package main

import (
	"log"
	"world-backup/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
