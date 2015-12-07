package main

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s -[c|p|r|e|l]\n", os.Args[0])
		return
	}
	client, _ := docker.NewClient(endpoint)
	switch os.Args[1] {
	case "-c":
		benchmarkVariesContainerNumber(client)
	case "-p":
		benchmarkVariesPeriod(client)
	case "-r":
		benchmarkVariesRoutineNumber(client)
	case "-e":
		benchmarkEventStream(client)
	case "-l":
		benchmarkEventLossRate(client)
	default:
	}
}
