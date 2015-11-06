package main

import (
	//"fmt"
	//"math"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

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
