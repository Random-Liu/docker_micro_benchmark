package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/random-liu/docker_micro_benchmark/event"
	"github.com/random-liu/docker_micro_benchmark/helpers"
)

var (
	wg = &sync.WaitGroup{}
)

func benchmarkVariesContainerNumber(client *docker.Client) {
	curDeadContainerNum := deadContainers[0]
	curAliveContainerNum := aliveContainers[0]
	helpers.CreateDeadContainers(client, curDeadContainerNum)
	helpers.CreateAliveContainers(client, curAliveContainerNum)
	for _, containerNum := range deadContainers {
		// Create more dead containers
		helpers.CreateDeadContainers(client, containerNum-curDeadContainerNum)
		curDeadContainerNum = containerNum
		// Get newest container ids
		containerIds := helpers.GetContainerIds(client)
		helpers.LogTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, true))
		helpers.LogLatency(helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true))
		helpers.LogTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, false))
		helpers.LogLatency(helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false))
		helpers.LogTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum))
		helpers.LogLatency(helpers.DoInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, containerIds))
	}

	for _, containerNum := range aliveContainers {
		// Create more alive containers
		helpers.CreateAliveContainers(client, containerNum-curAliveContainerNum)
		curAliveContainerNum = containerNum
		// Get newest container ids
		containerIds := helpers.GetContainerIds(client)
		helpers.LogTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, true))
		helpers.LogLatency(helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true))
		helpers.LogTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, false))
		helpers.LogLatency(helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false))
		helpers.LogTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum))
		helpers.LogLatency(helpers.DoInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, containerIds))
	}
}

func benchmarkVariesPeriod(client *docker.Client) {
	curAliveContainerNum := helpers.GetContainerNum(client, false)
	curDeadContainerNum := helpers.GetContainerNum(client, true) - curAliveContainerNum
	containerIds := helpers.GetContainerIds(client)
	for _, curPeriod := range listPeriods {
		helpers.LogTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			curPeriod, curDeadContainerNum, curAliveContainerNum, true))
		helpers.LogLatency(helpers.DoListContainerBenchMark(client, curPeriod, longTestPeriod, true))
		helpers.LogTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			curPeriod, curDeadContainerNum, curAliveContainerNum, false))
		helpers.LogLatency(helpers.DoListContainerBenchMark(client, curPeriod, longTestPeriod, false))
	}

	for _, curPeriod := range inspectPeriods {
		helpers.LogTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d]",
			curPeriod, curDeadContainerNum, curAliveContainerNum))
		helpers.LogLatency(helpers.DoInspectContainerBenchMark(client, curPeriod, shortTestPeriod, containerIds))
	}
}

func benchmarkVariesRoutineNumber(client *docker.Client) {
	curAliveContainerNum := helpers.GetContainerNum(client, false)
	curDeadContainerNum := helpers.GetContainerNum(client, true) - curAliveContainerNum
	containerIds := helpers.GetContainerIds(client)
	for _, curRoutineNumber := range routines {
		helpers.LogTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, No.Routines=%d, All=%v]",
			resyncPeriod, curDeadContainerNum, curAliveContainerNum, curRoutineNumber, true))
		helpers.LogLatency(helpers.DoParalListContainerBenchMark(client, resyncPeriod, shortTestPeriod, curRoutineNumber, true))
	}

	for _, curRoutineNumber := range routines {
		helpers.LogTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, No.Routines=%d]",
			routineInspectPeriod, curDeadContainerNum, curAliveContainerNum, curRoutineNumber))
		helpers.LogLatency(helpers.DoParalInspectContainerBenchMark(client, routineInspectPeriod, shortTestPeriod, curRoutineNumber, containerIds))
	}
}

func benchmarkEventStream(client *docker.Client) {
	for i, frequency := range eventFrequency {
		routineNumber := eventRoutines[i]
		helpers.LogTime(fmt.Sprintf("Event Stream Benchmark[Frequency=%v, No.Routines=%d]",
			frequency, routineNumber))
		var latencies []int
		stopchan := make(chan int, 1)
		defer close(stopchan)
		wg.Add(1)
		go func() {
			latencies, _ = helpers.DoEventStreamBenchMark(stopchan, client)
			wg.Done()
		}()
		cmd := exec.Command("event_generator/event_generator", strconv.Itoa(frequency), strconv.Itoa(routineNumber),
			strconv.FormatInt(shortTestPeriod.Nanoseconds(), 10), endpoint)
		if out, err := cmd.Output(); err != nil {
			panic(fmt.Sprintf("Error get output: %v", err))
		} else {
			cmd.Run()
			fmt.Print(string(out))
		}
		// Just make sure that all the events are received
		time.Sleep(time.Second)
		stopchan <- 1
		wg.Wait()
		helpers.LogLatency(latencies)
		helpers.LogTime(fmt.Sprintf("Event Stream Benchmark[Event Number=%v, Event Received=%v]", len(latencies)))
	}
}

func benchmarkEventLossRate(client *docker.Client) {
	for _, testPeriod := range testPeriodList {
		for t := 1; t <= timesForEachPeriod; t++ {
			helpers.LogTime(fmt.Sprintf("Event Stream Loss Rate Benchmark[Frequency=%v, No.Routines=%d, Period=%v, Times=%v]",
				defaultEventFrequency, defaultEventRoutines, testPeriod, t))
			var events []*docker.APIEvents
			stopchan := make(chan int, 1)
			defer close(stopchan)
			wg.Add(1)
			go func() {
				_, events = helpers.DoEventStreamBenchMark(stopchan, client)
				wg.Done()
			}()
			dockerIDs := event.StartGeneratingEvent(client, int64(defaultEventFrequency), defaultEventRoutines, testPeriod)
			// Just make sure that all the events are received
			time.Sleep(time.Second)
			stopchan <- 1
			wg.Wait()

			mismatchEventNum := 0
			dockerIDMap := map[string]bool{}
			for _, event := range events {
				if created, found := dockerIDMap[event.ID]; !found {
					if event.Status == "create" {
						dockerIDMap[event.ID] = true
					} else {
						helpers.LogTime(fmt.Sprintf("Detail Info[Mismatched Event: %+v", event))
						mismatchEventNum++
					}
				} else {
					if created && event.Status == "destroy" {
						dockerIDMap[event.ID] = false
					} else {
						helpers.LogTime(fmt.Sprintf("Detail Info[Mismatched Event: %+v", event))
						mismatchEventNum++
					}
				}
			}

			badOrderEventNum := 0
			var lastEvent *docker.APIEvents
			for _, event := range events {
				if lastEvent != nil && event.Time < lastEvent.Time {
					badOrderEventNum++
				}
				lastEvent = event
			}

			errDockerIDNum := 0
			for _, dockerID := range dockerIDs {
				if _, found := dockerIDMap[dockerID]; !found {
					errDockerIDNum++
				}
			}
			helpers.LogTime(fmt.Sprintf("Event Stream Loss Rate Benchmark[Event Created=%v, Event Received=%v, No.MisMatched Events=%v, No.Bad Order Events=%v, No.Error Docker IDs=%v]",
				len(dockerIDs)*2, len(events), mismatchEventNum, badOrderEventNum, errDockerIDNum))
		}
	}
}
