package main

import (
	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	usage := func() {
		fmt.Printf("Usage: %s -[o|c|i|r]\n", os.Args[0])
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
	case "-i":
		benchmarkVariesInterval(client)
	case "-r":
		benchmarkVariesRoutineNumber(client)
	default:
		usage()
	}
}
