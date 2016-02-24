package main

import (
	//	"fmt"
	//"os/exec"
	//"strconv"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	//"github.com/random-liu/docker_micro_benchmark/event"
	"github.com/random-liu/docker_micro_benchmark/helpers"
)

var (
	wg = &sync.WaitGroup{}
)

func benchmarkContainerStart(client *docker.Client) {
	cfg := containerStartConfig
	helpers.LogTitle("container_start")
	helpers.LogEVar(map[string]interface{}{"period": longTestPeriod})
	helpers.LogLabels("qps", "cps")
	for _, q := range cfg["qps"].([]float64) {
		start := time.Now()
		latencies := helpers.DoParalContainerStartBenchMark(client, q, longTestPeriod, cfg["routine"].(int))
		cps := float64(len(latencies)) / time.Now().Sub(start).Seconds()
		helpers.LogLatencyNew(latencies, helpers.Ftoas(q, cps)...)

		start = time.Now()
		latencies = helpers.DoParalContainerStopBenchMark(client, q, cfg["routine"].(int))
		cps = float64(len(latencies)) / time.Now().Sub(start).Seconds()
		helpers.LogLatencyNew(latencies, helpers.Ftoas(q, cps)...)
	}
}

func benchmarkVariesContainerNumber(client *docker.Client) {
	dead := deadContainers[0]
	alive := aliveContainers[0]
	ids := helpers.CreateDeadContainers(client, dead)
	ids = append(ids, helpers.CreateAliveContainers(client, alive)...)
	helpers.LogTitle("varies_container")
	helpers.LogEVar(map[string]interface{}{
		"period":   shortTestPeriod,
		"interval": defaultPeriod,
	})
	helpers.LogLabels("#dead", "#alive", "#total")
	for _, num := range deadContainers {
		// Create more dead containers
		ids = append(ids, helpers.CreateDeadContainers(client, num-dead)...)
		dead = num
		total := dead + alive
		latencies := helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true, nil)
		helpers.LogLatencyNew(latencies, helpers.Itoas(dead, alive, total)...)
		latencies = helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false, nil)
		helpers.LogLatencyNew(latencies, helpers.Itoas(dead, alive, total)...)
		latencies = helpers.DoInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, ids)
		helpers.LogLatencyNew(latencies, helpers.Itoas(dead, alive, total)...)
	}

	for _, num := range aliveContainers {
		// Create more alive containers
		ids = append(ids, helpers.CreateAliveContainers(client, num-alive)...)
		alive = num
		total := dead + alive
		latencies := helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true, nil)
		helpers.LogLatencyNew(latencies, helpers.Itoas(dead, alive, total)...)
		latencies = helpers.DoListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false, nil)
		helpers.LogLatencyNew(latencies, helpers.Itoas(dead, alive, total)...)
		latencies = helpers.DoInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, ids)
		helpers.LogLatencyNew(latencies, helpers.Itoas(dead, alive, total)...)
	}
}

func benchmarkVariesPeriod(client *docker.Client) {
	alive := helpers.GetContainerNum(client, false)
	dead := helpers.GetContainerNum(client, true) - alive
	containerIds := helpers.GetContainerIds(client)
	helpers.LogTitle("list_all")
	helpers.LogEVar(map[string]interface{}{
		"#alive": alive,
		"#dead":  dead,
		"all":    true,
		"period": longTestPeriod,
	})
	helpers.LogLabels("interval")
	for _, curPeriod := range listPeriods {
		latencies := helpers.DoListContainerBenchMark(client, curPeriod, longTestPeriod, true, nil)
		helpers.LogLatencyNew(latencies, helpers.Itoas(int(curPeriod/time.Millisecond))...)
	}

	helpers.LogTitle("list_alive")
	helpers.LogEVar(map[string]interface{}{
		"#alive": alive,
		"#dead":  dead,
		"all":    false,
		"period": longTestPeriod,
	})
	helpers.LogLabels("interval")
	for _, curPeriod := range listPeriods {
		latencies := helpers.DoListContainerBenchMark(client, curPeriod, longTestPeriod, false, nil)
		helpers.LogLatencyNew(latencies, helpers.Itoas(int(curPeriod/time.Millisecond))...)
	}

	helpers.LogTitle("inspect")
	helpers.LogEVar(map[string]interface{}{
		"#alive": alive,
		"#dead":  dead,
		"period": shortTestPeriod,
	})
	helpers.LogLabels("interval")
	for _, curPeriod := range inspectPeriods {
		latencies := helpers.DoInspectContainerBenchMark(client, curPeriod, shortTestPeriod, containerIds)
		helpers.LogLatencyNew(latencies, helpers.Itoas(int(curPeriod/time.Millisecond))...)
	}
}

func benchmarkVariesRoutineNumber(client *docker.Client) {
	alive := helpers.GetContainerNum(client, false)
	dead := helpers.GetContainerNum(client, true) - alive
	containerIds := helpers.GetContainerIds(client)

	helpers.LogTitle("list_all")
	helpers.LogEVar(map[string]interface{}{
		"#alive":           alive,
		"#dead":            dead,
		"all":              true,
		"interval":         resyncPeriod,
		"inspect-interval": routineInspectPeriod,
		"period":           shortTestPeriod,
	})
	helpers.LogLabels("#routines")
	for _, curRoutineNumber := range routines {
		latencies := helpers.DoParalListContainerBenchMark(client, resyncPeriod, shortTestPeriod, curRoutineNumber, true)
		helpers.LogLatencyNew(latencies, helpers.Itoas(curRoutineNumber)...)
	}

	helpers.LogTitle("inspect")
	helpers.LogEVar(map[string]interface{}{
		"#alive":   alive,
		"#dead":    dead,
		"interval": routineInspectPeriod,
		"period":   shortTestPeriod,
	})
	helpers.LogLabels("#routines")
	for _, curRoutineNumber := range routines {
		latencies := helpers.DoParalInspectContainerBenchMark(client, routineInspectPeriod, shortTestPeriod, curRoutineNumber, containerIds)
		helpers.LogLatencyNew(latencies, helpers.Itoas(curRoutineNumber)...)
	}
}

func benchmarkEventStream(client *docker.Client) {
	/*alive := helpers.GetContainerNum(client, false)
	dead := helpers.GetContainerNum(client, true) - alive

	// Benchmark Event Listener
	helpers.LogTitle("events")
	helpers.LogEVar(map[string]interface{}{
		"#alive": alive,
		"#dead":  dead,
		//"all":              true,
		//"interval":         resyncPeriod,
		//"inspect-interval": routineInspectPeriod,
		"period": shortTestPeriod,
	})

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
	}*/
}

func benchmarkEventLossRate(client *docker.Client) {
	/*for _, testPeriod := range testPeriodList {
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
	}*/
}
