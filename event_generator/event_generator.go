package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/random-liu/docker_micro_benchmark/helpers"
)

var (
	wg = &sync.WaitGroup{}
)

func StartEventGenerator(client *docker.Client, frequency int64, routineNumber int, testPeriod time.Duration) {
	period := time.Duration(time.Second.Nanoseconds() / frequency * int64(routineNumber))
	times := make([]int, routineNumber)
	startTime := time.Now()
	wg.Add(routineNumber)
	helpers.LogTime(fmt.Sprintf("Start Generating Event[Frequency=%v]", frequency))
	for id := 0; id < routineNumber; id++ {
		go func(id int) {
			client, _ = docker.NewClient("unix:///var/run/docker.sock")
			for {
				helpers.CreateAndRemoveContainers(client)
				times[id]++
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
	for _, time := range times {
		totalTimes += time
	}
	helpers.LogTime(fmt.Sprintf("Stop Generating Event[Expected Frequency=%v, Real Frequency=%v]", frequency, float64(totalTimes)/testPeriod.Seconds()))
}

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s Frequency Routines TestPeriod EndPoint\n", os.Args[0])
		return
	}
	var frequency, testPeriod int64
	var routineNumber int
	var client *docker.Client
	var err error
	frequency, err = strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Error get frequency: %v", err))
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
	StartEventGenerator(client, frequency, routineNumber, time.Duration(testPeriod))
}
