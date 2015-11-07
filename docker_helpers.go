package main

import (
	"math/rand"
	"strconv"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func createDeadContainers(client *docker.Client, num int) {
	for i := 0; i < num; i++ {
		name := "benchmark_container_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		dockerOpts := docker.CreateContainerOptions{
			Name: name,
			Config: &docker.Config{
				Image: "ubuntu",
			},
		}
		container, err := client.CreateContainer(dockerOpts)
		if err != nil {
			panic("Error create containers")
		}
		client.StartContainer(container.ID, &docker.HostConfig{})
	}
}

func createAliveContainers(client *docker.Client, num int) {
	for i := 0; i < num; i++ {
		name := "benchmark_container_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		dockerOpts := docker.CreateContainerOptions{
			Name: name,
			Config: &docker.Config{
				AttachStderr: false,
				AttachStdin:  false,
				AttachStdout: false,
				Tty:          true,
				Cmd:          []string{"/bin/bash"},
				Image:        "ubuntu",
			},
		}
		container, err := client.CreateContainer(dockerOpts)
		if err != nil {
			panic("Error create containers")
		}
		client.StartContainer(container.ID, &docker.HostConfig{})
	}
}

func doListContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, all bool) {
	startTime := time.Now()
	latencies := []int{}
	for {
		start := time.Now()
		client.ListContainers(docker.ListContainersOptions{All: all})
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

func doInspectContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration) {
	containers, _ := client.ListContainers(docker.ListContainersOptions{All: true})
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

/*
func doListContainers(client *docker.Client, all bool, sync bool) {
	for _, curPeriod := range periods {
		logTime(fmt.Sprintf("Test starts with period=%v", curPeriod))
		var testTimes int
		maxLatency := 0
		minLatency := math.MaxInt64
		totalLatency := 0
		startTime := time.Now()
		for {
			testTimes++
			start := time.Now()
			client.ListContainers(docker.ListContainersOptions{All: all})
			end := time.Now()
			curLatency := int(end.Sub(start).Nanoseconds())
			if maxLatency < curLatency {
				maxLatency = curLatency
			}
			if minLatency > curLatency {
				minLatency = curLatency
			}
			totalLatency += curLatency
			if time.Now().Sub(startTime) >= testPeriod {
				break
			}
			if curPeriod != 0 {
				time.Sleep(curPeriod)
			}
		}
		logLatency(totalLatency/testTimes, maxLatency, minLatency)
	}
	if sync {
		wg.Done()
	}
}

func benchmarkListContainers(client *docker.Client, routineNumber int, all bool) {
	var adj string
	if all {
		adj = "All"
	} else {
		adj = "Alive"
	}
	logTime(fmt.Sprintf("Test List %s Containers with routineNumber %d", adj, routineNumber))
	if routineNumber == 0 {
		doListContainers(client, all, false)
	} else {
		wg.Add(routineNumber)
		for i := 0; i < routineNumber; i++ {
			go doListContainers(client, all, true)
		}
		wg.Wait()
	}
}

func listBenchmarkPlan(client *docker.Client) {
	// Test listContainers (Like PLEG)
	benchmarkListContainers(client, 0, true)
	benchmarkListContainers(client, 0, false)

	// Test listContainers in each mutilple routines (Like PodWorkers)
	for _, routineNumber := range routines {
		// PodWorkers only use list containers true
		benchmarkListContainers(client, routineNumber, true)
	}
}
*/
