package main

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

func benchmarkVariesContainerNumber(client *docker.Client) {
	curDeadContainerNum := deadContainers[0]
	curAliveContainerNum := aliveContainers[0]
	createDeadContainers(client, curDeadContainerNum)
	createAliveContainers(client, curAliveContainerNum)
	for _, containerNum := range deadContainers {
		// Create more dead containers
		createDeadContainers(client, containerNum-curDeadContainerNum)
		curDeadContainerNum = containerNum
		// Get newest container ids
		containerIds := getContainerIds(client)
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, true))
		logLatency(doListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true))
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, false))
		logLatency(doListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false))
		logTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum))
		logLatency(doInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, containerIds))
	}

	for _, containerNum := range aliveContainers {
		// Create more alive containers
		createAliveContainers(client, containerNum-curAliveContainerNum)
		curAliveContainerNum = containerNum
		// Get newest container ids
		containerIds := getContainerIds(client)
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, true))
		logLatency(doListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true))
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, false))
		logLatency(doListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false))
		logTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum))
		logLatency(doInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, containerIds))
	}
}

func benchmarkVariesPeriod(client *docker.Client) {
	curDeadContainerNum := deadContainers[len(deadContainers)-1]
	curAliveContainerNum := aliveContainers[len(aliveContainers)-1]
	containerIds := getContainerIds(client)
	for _, curPeriod := range listPeriods {
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			curPeriod, curDeadContainerNum, curAliveContainerNum, true))
		logLatency(doListContainerBenchMark(client, curPeriod, longTestPeriod, true))
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			curPeriod, curDeadContainerNum, curAliveContainerNum, false))
		logLatency(doListContainerBenchMark(client, curPeriod, longTestPeriod, false))
	}

	for _, curPeriod := range inspectPeriods {
		logTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d]",
			curPeriod, curDeadContainerNum, curAliveContainerNum))
		logLatency(doInspectContainerBenchMark(client, curPeriod, shortTestPeriod, containerIds))
	}
}

func benchmarkVariesRoutineNumber(client *docker.Client) {
	curDeadContainerNum := deadContainers[len(deadContainers)-1]
	curAliveContainerNum := aliveContainers[len(aliveContainers)-1]
	containerIds := getContainerIds(client)
	fmt.Println(containerIds)
	for _, curRoutineNumber := range routines {
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, No.Routines=%d, All=%v]",
			resyncPeriod, curDeadContainerNum, curAliveContainerNum, curRoutineNumber, true))
		logLatency(doParalListContainerBenchMark(client, resyncPeriod, shortTestPeriod, curRoutineNumber, true))
	}

	for _, curRoutineNumber := range routines {
		logTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, No.Routines=%d]",
			routineInspectPeriod, curDeadContainerNum, curAliveContainerNum, curRoutineNumber))
		logLatency(doParalInspectContainerBenchMark(client, routineInspectPeriod, shortTestPeriod, curRoutineNumber, containerIds))
	}
}
