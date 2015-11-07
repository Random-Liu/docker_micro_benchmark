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
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, true))
		doListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true)
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, false))
		doListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false)
		logTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum))
		doInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod)
	}

	for _, containerNum := range aliveContainers {
		// Create more alive containers
		createAliveContainers(client, containerNum-curAliveContainerNum)
		curAliveContainerNum = containerNum
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, true))
		doListContainerBenchMark(client, defaultPeriod, shortTestPeriod, true)
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, false))
		doListContainerBenchMark(client, defaultPeriod, shortTestPeriod, false)
		logTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum))
		doInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod)
	}
}

func benchmarkVariesPeriod(client *docker.Client) {
	curDeadContainerNum := deadContainers[len(deadContainers)-1]
	curAliveContainerNum := aliveContainers[len(aliveContainers)-1]
	for _, curPeriod := range periods {
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			curPeriod, curDeadContainerNum, curAliveContainerNum, true))
		doListContainerBenchMark(client, defaultPeriod, longTestPeriod, true)
		logTime(fmt.Sprintf("List Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d, All=%v]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum, false))
		doListContainerBenchMark(client, defaultPeriod, longTestPeriod, false)
		logTime(fmt.Sprintf("Inspect Benchmark[Period=%v, No.DeadContainers=%d, No.AliveContainers=%d]",
			defaultPeriod, curDeadContainerNum, curAliveContainerNum))
		doInspectContainerBenchMark(client, defaultPeriod, longTestPeriod)
	}
}
