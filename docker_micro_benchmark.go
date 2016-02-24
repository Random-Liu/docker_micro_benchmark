package main

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"os"
)

func main() {
	usage := func() {
		fmt.Printf("Usage: %s -[o|c|p|r|e|l]\n", os.Args[0])
	}
	if len(os.Args) != 2 {
		usage()
		return
	}
	client, _ := docker.NewClient(endpoint)
	client.PullImage(docker.PullImageOptions{Repository: "ubuntu", Tag: "latest"}, docker.AuthConfiguration{})
	switch os.Args[1] {
	case "-o":
		benchmarkContainerStart(client)
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
		usage()
	}
}
