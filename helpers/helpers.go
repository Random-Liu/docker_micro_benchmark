package helpers

import (
	"fmt"
	"math"
	"sort"
	"time"
)

var rates = []float64{0.5, 0.75, 0.95, 0.99}

var last = time.Now()

func LogTime(label string) {
	now := time.Now()
	fmt.Printf("%02d:%02d:%02d\t%s\n", last.Hour(), last.Minute(), last.Second(), label)
	fmt.Printf("%02d:%02d:%02d\t%s\n", now.Hour(), now.Minute(), now.Second(), label)
	last = now
}

func LogEVar(vars map[string]interface{}) {
	for k, v := range vars {
		fmt.Printf("%s=%v ", k, v)
	}
	fmt.Println()
}

func LogLabels(labels string) {
	fmt.Printf("time\t%%50\t%%75\t%%95\t%%99\t%s\n", labels)
}

func LogLatencyNew(variables string, latencies []int) {
	if len(latencies) <= 0 {
		panic("No latency record!")
	}
	sort.Ints(latencies)
	average := func(latencies []int) int {
		total := 0
		for _, l := range latencies {
			total += l
		}
		return total / len(latencies)
	}

	var avgs [4]float64
	for i, rate := range rates {
		n := int(math.Ceil((1 - rate) * float64(len(latencies))))
		avgs[i] = float64(average(latencies[len(latencies)-n:])) / 1000000
	}
	LogTime(fmt.Sprintf("%.2f\t%.2f\t%.2f\t%.2f\t%s", avgs[0], avgs[1], avgs[2], avgs[3], variables))
}

func LogLatency(latencies []int) {
	if len(latencies) <= 0 {
		panic("No latency record!")
	}
	sort.Ints(latencies)
	average := func(latencies []int) int {
		total := 0
		for _, l := range latencies {
			total += l
		}
		return total / len(latencies)
	}

	var avgs [4]float64
	for i, rate := range rates {
		n := int(math.Ceil((1 - rate) * float64(len(latencies))))
		avgs[i] = float64(average(latencies[len(latencies)-n:])) / 1000000
	}
	LogTime(fmt.Sprintf("%f %f %f %f", avgs[0], avgs[1], avgs[2], avgs[3]))
}
