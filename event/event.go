package event

import (
	"fmt"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/random-liu/docker_micro_benchmark/helpers"
)

var (
	wg = &sync.WaitGroup{}
)

func StartGeneratingEvent(client *docker.Client, frequency int64, routineNumber int, testPeriod time.Duration) []string {
	period := time.Duration(time.Second.Nanoseconds() / frequency * int64(routineNumber))
	groupedDockerIDs := make([][]string, routineNumber)
	startTime := time.Now()
	wg.Add(routineNumber)
	helpers.LogTime(fmt.Sprintf("Start Generating Event[Frequency=%v]", frequency))
	for id := 0; id < routineNumber; id++ {
		go func(id int) {
			client, _ = docker.NewClient("unix:///var/run/docker.sock")
			for {
				dockerID := helpers.CreateAndRemoveContainers(client)
				groupedDockerIDs[id] = append(groupedDockerIDs[id], dockerID)
				if time.Now().Sub(startTime) >= testPeriod {
					break
				}
				time.Sleep(period)
			}
			wg.Done()
		}(id)
	}
	wg.Wait()
	totalTimes := 0
	dockerIDs := []string{}
	for _, group := range groupedDockerIDs {
		dockerIDs = append(dockerIDs, group...)
		totalTimes += len(group)
	}
	helpers.LogTime(fmt.Sprintf("Stop Generating Event[Expected Frequency=%v, Real Frequency=%v, Total Event Number=%v]", frequency, float64(totalTimes)/testPeriod.Seconds(), totalTimes*2))
	return dockerIDs
}
