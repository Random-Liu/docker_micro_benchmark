package main

import (
	//	"fmt"
	/*	"math"
		"sync"
		"time"
	*/

	docker "github.com/fsouza/go-dockerclient"
)

/*
var (
	before map[string]int64 = map[string]int64{}
	after  map[string]int64 = map[string]int64{}
	times  map[string]int64 = map[string]int64{}
)
*/

// TODO Use configuration file later
/*const (
	minPeriod  = 10 * time.Second
	maxPeriod  = 10 * time.Second
	accPeriod  = time.Second
	testPeriod = 100 * time.Second
)

var (
	wg = &sync.WaitGroup{}
)

func logTime(label string) {
	now := time.Now()
	fmt.Printf("%d:%d:%d\t%s\n", now.Hour(), now.Minute(), now.Second(), label)
}

func logLatency(averageLatency, maxLatency, minLatency int) {
	fmt.Printf("averageLatency:%dms, maxLatency:%dms, minLatency:%dms\n", averageLatency/1000000, maxLatency/1000000, minLatency/1000000)
}

func testListContainers(client *docker.Client, curPeriod time.Duration, all bool) {
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
	logTime(fmt.Sprintf("Test ends with period=%v", curPeriod))
	wg.Done()
}
*/
func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	benchmarkVariesContainerNumber(client)
}
