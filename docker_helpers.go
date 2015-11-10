package main

import (
	"fmt"
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
			panic(fmt.Sprintf("Error create containers: %v", err))
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
			panic(fmt.Sprintf("Error create containers: %v", err))
		}
		client.StartContainer(container.ID, &docker.HostConfig{})
	}
}

func doListContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, all bool) []int {
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
	return latencies
}

func doInspectContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, containerIds []string) []int {
	startTime := time.Now()
	latencies := []int{}
	rand.Seed(time.Now().Unix())
	for {
		containerId := containerIds[rand.Int()%len(containerIds)]
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
	return latencies
}

// Use true because that's the behaviour of the pod worker
func doParalListContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, routineNumber int, all bool) []int {
	wg.Add(routineNumber)
	latenciesTable := make([][]int, routineNumber)
	for i := 0; i < routineNumber; i++ {
		go func(index int) {
			latenciesTable[index] = doListContainerBenchMark(client, curPeriod, testPeriod, all)
			wg.Done()
		}(i)
	}
	wg.Wait()
	allLatencies := []int{}
	for _, latencies := range latenciesTable {
		allLatencies = append(allLatencies, latencies...)
	}
	return allLatencies
}

func doParalInspectContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, routineNumber int, containerIds []string) []int {
	wg.Add(routineNumber)
	latenciesTable := make([][]int, routineNumber)
	for i := 0; i < routineNumber; i++ {
		go func(index int) {
			latenciesTable[index] = doInspectContainerBenchMark(client, curPeriod, testPeriod, containerIds)
			wg.Done()
		}(i)
	}
	wg.Wait()
	allLatencies := []int{}
	for latencies := range latenciesTable {
		allLatencies = append(allLatencies, latencies)
	}
	return allLatencies
}

func getContainerIds(client *docker.Client) (containerIds []string) {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		panic(fmt.Sprintf("Error list containers: %v", err))
	}
	for _, container := range containers {
		containerIds = append(containerIds, container.ID)
	}
	return containerIds
}
