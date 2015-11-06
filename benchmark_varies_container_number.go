package main

import (
	"fmt"
	"strconv"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func createDeadContainers(client *docker.Client, num int) {
	for i := 0; i < num; i++ {
		name := "benchmark_container_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		dockerOpts := docker.CreateContainerOptions{
			Name: name,
			Config: &docker.Config{
				Image: "ubuntu",
			},
		}
		container, err := client.CreateContainer(dockerOpts)
		if err != nil {
			panic("Error create containers")
		}
		client.StartContainer(container.ID, &docker.HostConfig{})
	}
}

func createAliveContainers(client *docker.Client, num int) {
	for i := 0; i < num; i++ {
		name := "benchmark_container_" + strconv.FormatInt(time.Now().UnixNano(), 10)
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
			panic("Error create containers")
		}
		client.StartContainer(container.ID, &docker.HostConfig{})
	}
}

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
		doInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, false)
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
		doInspectContainerBenchMark(client, defaultPeriod, shortTestPeriod, false)
	}
}
