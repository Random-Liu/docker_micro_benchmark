package helpers

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	wg = &sync.WaitGroup{}
)

func newContainerName() string {
	return "benchmark_container_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func CreateAndRemoveContainers(client *docker.Client) {
	name := newContainerName()
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
	removeOpts := docker.RemoveContainerOptions{
		ID: container.ID,
	}
	if err := client.RemoveContainer(removeOpts); err != nil {
		panic(fmt.Sprintf("Error remove containers: %v", err))
	}
}

func CreateDeadContainers(client *docker.Client, num int) {
	for i := 0; i < num; i++ {
		name := newContainerName()
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

func CreateAliveContainers(client *docker.Client, num int) {
	for i := 0; i < num; i++ {
		name := newContainerName()
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

func DoListContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, all bool) []int {
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

func DoInspectContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, containerIds []string) []int {
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
func DoParalListContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, routineNumber int, all bool) []int {
	wg.Add(routineNumber)
	latenciesTable := make([][]int, routineNumber)
	for i := 0; i < routineNumber; i++ {
		go func(index int) {
			latenciesTable[index] = DoListContainerBenchMark(client, curPeriod, testPeriod, all)
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

func DoParalInspectContainerBenchMark(client *docker.Client, curPeriod, testPeriod time.Duration, routineNumber int, containerIds []string) []int {
	wg.Add(routineNumber)
	latenciesTable := make([][]int, routineNumber)
	for i := 0; i < routineNumber; i++ {
		go func(index int) {
			latenciesTable[index] = DoInspectContainerBenchMark(client, curPeriod, testPeriod, containerIds)
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

func DoEventStreamBenchMark(stopchan chan int, client *docker.Client) []int {
	eventchan := make(chan *docker.APIEvents, 1000)
	defer close(eventchan)
	if err := client.AddEventListener(eventchan); err != nil {
		panic(fmt.Sprintf("Error add event listener: %v", err))
	}
	latencies := []int{}
	for {
		select {
		case event := <-eventchan:
			latency := time.Now().Unix() - event.Time
			latencies = append(latencies, int(latency))
		case <-stopchan:
			client.RemoveEventListener(eventchan)
			return latencies
		}
	}
}

func GetContainerIds(client *docker.Client) (containerIds []string) {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		panic(fmt.Sprintf("Error list containers: %v", err))
	}
	for _, container := range containers {
		containerIds = append(containerIds, container.ID)
	}
	return containerIds
}

func GetContainerNum(client *docker.Client, all bool) int {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: all})
	if err != nil {
		panic(fmt.Sprintf("Error list containers: %v", err))
	}
	return len(containers)
}