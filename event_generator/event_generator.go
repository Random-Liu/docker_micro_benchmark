package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/random-liu/docker_micro_benchmark/event"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s Interval Routines TestPeriod EndPoint\n", os.Args[0])
		return
	}
	var interval, testPeriod int64
	var routineNumber int
	var client *docker.Client
	var err error
	interval, err = strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Error get interval: %v", err))
	}
	routineNumber, err = strconv.Atoi(os.Args[2])
	if err != nil {
		panic(fmt.Sprintf("Error get routine number: %v", err))
	}
	testPeriod, err = strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Error get test period: %v", err))
	}
	client, err = docker.NewClient(os.Args[4])
	if err != nil {
		panic(fmt.Sprintf("Error create docker client: %v", err))
	}
	event.StartGeneratingEvent(client, time.Duration(interval), routineNumber, time.Duration(testPeriod))
}
