package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg = &sync.WaitGroup{}
)

func logTime(label string) {
	now := time.Now()
	fmt.Printf("%02d:%02d:%02d\t%s\n", now.Hour(), now.Minute(), now.Second(), label)
}

func logLatency(latencies []int) {
	if len(latencies) < 0 {
		panic("No latency record!")
	}
	maxLatency := latencies[0]
	minLatency := latencies[0]
	totalLatency := 0
	for _, latency := range latencies {
		if maxLatency < latency {
			maxLatency = latency
		}
		if minLatency > latency {
			minLatency = latency
		}
		totalLatency += latency
	}
	averageLatency := totalLatency / len(latencies)
	logTime(fmt.Sprintf("averageLatency:%fms, maxLatency:%fms, minLatency:%fms",
		float64(averageLatency)/1000000, float64(maxLatency)/1000000, float64(minLatency)/1000000))
}
