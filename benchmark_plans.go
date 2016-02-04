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
	eVars := map[string]interface{}{
		"period":   shortTestPeriod,
		"interval": defaultPeriod,
	}
	helpers.LogEVar(eVars)
	helpers.LogLabels("#dead\t#alive\t#total")
	for _, containerNum := range deadContainers {
		// Create more dead containers
		helpers.CreateDeadContainers(client, containerNum-curDeadContainerNum)
		curDeadContainerNum = containerNum
		// Get newest container ids
		containerIds := helpers.GetContainerIds(client)
		helpers.LogLatencyNew(fmt.Sprintf("%d\t%d\t%d", curDeadContainerNum, curAliveContainerNum, curDeadContainerNum+curAliveContainerNum), helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true, nil))
		helpers.LogLatencyNew(fmt.Sprintf("%d\t%d\t%d", curDeadContainerNum, curAliveContainerNum, curDeadContainerNum+curAliveContainerNum), helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false, nil))
		helpers.LogLatencyNew(fmt.Sprintf("%d\t%d\t%d", curDeadContainerNum, curAliveContainerNum, curDeadContainerNum+curAliveContainerNum), helpers.DoInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, containerIds))
	}

	for _, containerNum := range aliveContainers {
		// Create more alive containers
		helpers.CreateAliveContainers(client, containerNum-curAliveContainerNum)
		curAliveContainerNum = containerNum
		// Get newest container ids
		containerIds := helpers.GetContainerIds(client)
		helpers.LogLatencyNew(fmt.Sprintf("%d\t%d\t%d", curDeadContainerNum, curAliveContainerNum, curDeadContainerNum+curAliveContainerNum), helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true, nil))
		helpers.LogLatencyNew(fmt.Sprintf("%d\t%d\t%d", curDeadContainerNum, curAliveContainerNum, curDeadContainerNum+curAliveContainerNum), helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false, nil))
		helpers.LogLatencyNew(fmt.Sprintf("%d\t%d\t%d", curDeadContainerNum, curAliveContainerNum, curDeadContainerNum+curAliveContainerNum), helpers.DoInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, containerIds))
	}
}

func benchmarkVariesPeriod(client *docker.Client) {
	curAliveContainerNum := helpers.GetContainerNum(client, false)
	curDeadContainerNum := helpers.GetContainerNum(client, true) - curAliveContainerNum
	containerIds := helpers.GetContainerIds(client)
	eVars := map[string]interface{}{
		"#alive": curAliveContainerNum,
		"#dead":  curDeadContainerNum,
		"period": longTestPeriod,
	}
	helpers.LogEVar(eVars)
	helpers.LogLabels("interval")

	for _, curPeriod := range listPeriods {
		helpers.LogLatencyNew(fmt.Sprintf("%d", curPeriod/time.Millisecond), helpers.DoListContainerBenchMark(client, curPeriod, longTestPeriod, true, nil))
	}

	for _, curPeriod := range listPeriods {
		helpers.LogLatencyNew(fmt.Sprintf("%d", curPeriod/time.Millisecond), helpers.DoListContainerBenchMark(client, curPeriod, longTestPeriod, false, nil))
	}

	for _, curPeriod := range inspectPeriods {
		helpers.LogLatencyNew(fmt.Sprintf("%d", curPeriod/time.Millisecond), helpers.DoInspectContainerBenchMark(client, curPeriod, longTestPeriod, containerIds))
	}
}

func benchmarkVariesRoutineNumber(client *docker.Client) {
	curAliveContainerNum := helpers.GetContainerNum(client, false)
	curDeadContainerNum := helpers.GetContainerNum(client, true) - curAliveContainerNum
	containerIds := helpers.GetContainerIds(client)
	eVars := map[string]interface{}{
		"#alive":           curAliveContainerNum,
		"#dead":            curDeadContainerNum,
		"all":              true,
		"list interval":    resyncPeriod,
		"inspect interval": routineInspectPeriod,
		"period":           shortTestPeriod,
	}
	helpers.LogEVar(eVars)
	helpers.LogLabels("#routines")
	for _, curRoutineNumber := range routines {
		helpers.LogLatencyNew(fmt.Sprintf("%d", curRoutineNumber), helpers.DoParalListContainerBenchMark(client, resyncPeriod, shortTestPeriod, curRoutineNumber, true))
	}
	for _, curRoutineNumber := range routines {
		helpers.LogLatencyNew(fmt.Sprintf("%d", curRoutineNumber), helpers.DoParalInspectContainerBenchMark(client, routineInspectPeriod, shortTestPeriod, curRoutineNumber, containerIds))
	}
}

func benchmarkEventStream(client *docker.Client) {
	// Benchmark Event Listener
	helpers.LogTime("Event Stream Benchmark - Event Listener")
	for i, frequency := range eventFrequency {
		routineNumber := eventRoutines[i]
		helpers.LogTime(fmt.Sprintf("Event Stream Benchmark[Frequency=%v, No.Routines=%d, TestPeriod=%v]",
			frequency, routineNumber, shortTestPeriod))
		var latencies []int
		stopchan := make(chan int, 0)
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
		close(stopchan)
		wg.Wait()
		helpers.LogLatency(latencies)
		helpers.LogTime(fmt.Sprintf("Event Stream Benchmark[Event Number=%v]", len(latencies)))
	}

	// Benchmark Relist
	helpers.LogTime("Event Stream Benchmark - Relist")
	for i, frequency := range eventFrequency {
		routineNumber := eventRoutines[i]
		helpers.LogTime(fmt.Sprintf("Event Stream Benchmark[Frequency=%v, No.Routines=%d, ResyncPeriod=%v, TestPeriod=%v]",
			frequency, routineNumber, resyncPeriod, shortTestPeriod))
		stopchan := make(chan int, 0)
		wg.Add(1)
		go func() {
			_ = helpers.DoListContainerBenchMark(client, resyncPeriod, shortTestPeriod, true, stopchan)
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
		close(stopchan)
		wg.Wait()
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

			badOrderEventNum := 0
			var lastEvent *docker.APIEvents
			for _, event := range events {
				if lastEvent != nil && event.Time < lastEvent.Time {
					badOrderEventNum++
				}
				lastEvent = event
			}

			dockerIDMap := map[string][]*docker.APIEvents{}
			for _, event := range events {
				dockerIDMap[event.ID] = append(dockerIDMap[event.ID], event)
			}

			errorEventNum := 0
			errorDockerNum := 0     // Total Error Docker Num
			extraEventDocker := 0   // Docker with duplicated events
			missingEventDocker := 0 // Docker with missed events
			orderWrongTimeRightEventDocker := 0
			orderRightTimeWrongEventDocker := 0
			orderWrongTimeWrongEventDocker := 0
			rightEventDocker := 0 // Docker with right events
			for _, dockerID := range dockerIDs {
				eventsPerDocker := dockerIDMap[dockerID]
				if len(eventsPerDocker) < 2 {
					missingEventDocker++
				} else if len(eventsPerDocker) > 2 {
					extraEventDocker++
				} else {
					if eventsPerDocker[0].Status == "create" && eventsPerDocker[1].Status == "destroy" {
						if eventsPerDocker[0].Time <= eventsPerDocker[1].Time {
							rightEventDocker++
							continue
						} else {
							orderRightTimeWrongEventDocker++
						}
					} else if eventsPerDocker[0].Status == "destroy" && eventsPerDocker[1].Status == "create" {
						if eventsPerDocker[0].Time >= eventsPerDocker[1].Time {
							orderWrongTimeRightEventDocker++
						} else {
							orderWrongTimeWrongEventDocker++
						}
					}
				}
				errorEventNum += len(eventsPerDocker)
				errorDockerNum++
			}
			helpers.LogTime(fmt.Sprintf("Event Stream Loss Rate Benchmark[Event Created=%v, Event Received=%v, No.Error Events=%v, No.Error Docker=%v, No.Bad Order Events=%v, No.Extra Event Docker=%v, No.Missing Event Docker=%v, No.Right Order Wrong Time Event Docker=%v, No.Wrong Order Right Time Event Docker=%v, No.Wrong Order Wrong Time Event Docker=%v]",
				len(dockerIDs)*2, len(events), errorEventNum, errorDockerNum, badOrderEventNum, extraEventDocker, missingEventDocker, orderRightTimeWrongEventDocker, orderWrongTimeRightEventDocker, orderWrongTimeWrongEventDocker))
		}
	}
}
