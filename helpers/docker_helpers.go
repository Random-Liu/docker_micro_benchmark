package helpers

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/juju/ratelimit"
)

var (
	wg = &sync.WaitGroup{}
)

func newContainerName() string {
	return "benchmark_container_" + strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Int())
}

func CreateContainers(client *docker.Client, num int) []string {
	ids := []string{}
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
		ids = append(ids, container.ID)
	}
	return ids
}

func StartContainers(client *docker.Client, ids []string) {
	for _, id := range ids {
		client.StartContainer(id, &docker.HostConfig{})
	}
}

func StopContainers(client *docker.Client, ids []string) {
	for _, id := range ids {
		client.StopContainer(id, 10)
	}
}

func RemoveContainers(client *docker.Client, ids []string) {
	for _, id := range ids {
		removeOpts := docker.RemoveContainerOptions{
			ID: id,
		}
		if err := client.RemoveContainer(removeOpts); err != nil {
			panic(fmt.Sprintf("Error remove containers: %v", err))
		}
	}
}

func CreateAndRemoveContainers(client *docker.Client) string {
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
	return container.ID
}

func CreateDeadContainers(client *docker.Client, num int) []string {
	return CreateContainers(client, num)
}

func CreateAliveContainers(client *docker.Client, num int) []string {
	ids := CreateContainers(client, num)
	StartContainers(client, ids)
	return ids
}

func DoListContainerBenchMark(client *docker.Client, interval, testPeriod time.Duration, all bool, stopchan chan int) []int {
	startTime := time.Now()
	latencies := []int{}
	for {
		start := time.Now()
		client.ListContainers(docker.ListContainersOptions{All: all})
		end := time.Now()
		latencies = append(latencies, int(end.Sub(start).Nanoseconds()))
		if stopchan == nil {
			if time.Now().Sub(startTime) >= testPeriod {
				return latencies
			}
		} else {
			select {
			case <-stopchan:
				return latencies
			default:
			}
		}
		if interval != 0 {
			time.Sleep(interval)
		}
	}
	return latencies
}

func DoInspectContainerBenchMark(client *docker.Client, interval, testPeriod time.Duration, containerIds []string) []int {
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
		if interval != 0 {
			time.Sleep(interval)
		}
	}
	return latencies
}

// Use true because that's the behaviour of the pod worker
func DoParalListContainerBenchMark(client *docker.Client, interval, testPeriod time.Duration, routineNumber int, all bool) []int {
	wg.Add(routineNumber)
	latenciesTable := make([][]int, routineNumber)
	for i := 0; i < routineNumber; i++ {
		go func(index int) {
			latenciesTable[index] = DoListContainerBenchMark(client, interval, testPeriod, all, nil)
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

func DoParalInspectContainerBenchMark(client *docker.Client, interval, testPeriod time.Duration, routineNumber int, containerIds []string) []int {
	wg.Add(routineNumber)
	latenciesTable := make([][]int, routineNumber)
	for i := 0; i < routineNumber; i++ {
		go func(index int) {
			latenciesTable[index] = DoInspectContainerBenchMark(client, interval, testPeriod, containerIds)
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

func DoParalContainerStartBenchMark(client *docker.Client, qps float64, testPeriod time.Duration, routineNumber int) []int {
	wg.Add(routineNumber)
	ratelimit := ratelimit.NewBucketWithRate(qps, int64(routineNumber))
	latenciesTable := make([][]int, routineNumber)
	for i := 0; i < routineNumber; i++ {
		go func(index int) {
			startTime := time.Now()
			latencies := []int{}
			for {
				ratelimit.Wait(1)
				start := time.Now()
				ids := CreateContainers(client, 1)
				StartContainers(client, ids)
				end := time.Now()
				latencies = append(latencies, int(end.Sub(start).Nanoseconds()))
				if time.Now().Sub(startTime) >= testPeriod {
					break
				}
			}
			latenciesTable[index] = latencies
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

func DoParalContainerStopBenchMark(client *docker.Client, qps float64, routineNumber int) []int {
	ids := GetContainerIds(client)
	idTable := make([][]string, routineNumber)
	for i := 0; i < len(ids); i++ {
		idTable[i%routineNumber] = append(idTable[i%routineNumber], ids[i])
	}
	wg.Add(routineNumber)
	ratelimit := ratelimit.NewBucketWithRate(qps, int64(routineNumber))
	latenciesTable := make([][]int, routineNumber)
	for i := 0; i < routineNumber; i++ {
		go func(index int) {
			latencies := []int{}
			for _, id := range idTable[index] {
				ratelimit.Wait(1)
				start := time.Now()
				StopContainers(client, []string{id})
				RemoveContainers(client, []string{id})
				end := time.Now()
				latencies = append(latencies, int(end.Sub(start).Nanoseconds()))
			}
			latenciesTable[index] = latencies
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
