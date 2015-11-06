package main

import (
	"math/rand"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func doInspectContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, all bool) {
	containers, _ := client.ListContainers(docker.ListContainersOptions{All: all})
	containersId := []string{}
	for _, container := range containers {
		containersId = append(containersId, container.ID)
	}
	startTime := time.Now()
	latencies := []int{}
	rand.Seed(time.Now().Unix())
	for {
		containerId := containersId[rand.Int()%len(containersId)]
		start := time.Now()
		client.InspectContainer(containerId)
		end := time.Now()
		latencies = append(latencies, int(end.Sub(start).Nanoseconds()))
		if time.Now().Sub(startTime) >= testPeriod {
			break
		}
		if curPeriod != 0 {
			time.Sleep(curPeriod)
		}
	}
	logLatency(latencies)
}
