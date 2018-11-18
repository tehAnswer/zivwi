package main

import (
	"os"

	app "github.com/tehAnswer/zivwi"
	worker "github.com/tehAnswer/zivwi/worker"
)

func main() {
	concept := os.Args[1]
	if concept == "worker" {
		worker.Run()
	} else {
		app.Run()
	}
}
