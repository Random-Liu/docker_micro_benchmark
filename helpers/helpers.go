package helpers

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

var rates = []float64{0.5, 0.75, 0.95, 0.99}

var last = time.Now()

func LogTime(label string) {
	now := time.Now()
	fmt.Printf("%02d:%02d:%02d:%02d\t%s\n", last.Day(), last.Hour(), last.Minute(), last.Second(), label)
	fmt.Printf("%02d:%02d:%02d:%02d\t%s\n", now.Day(), now.Hour(), now.Minute(), now.Second(), label)
	last = now
}

func LogTitle(title string) {
	fmt.Println()
	fmt.Println(title)
}

func LogEVar(vars map[string]interface{}) {
	for k, v := range vars {
		fmt.Printf("%s=%v ", k, v)
	}
	fmt.Println()
}

func LogLabels(labels ...string) {
	fmt.Printf("time\t%%50\t%%75\t%%95\t%%99\t%s\n", strings.Join(labels, "\t"))
}

func LogLatencyNew(latencies []int, variables ...string) {
	average := func(latencies []int) int {
		if len(latencies) <= 0 {
			return 0
		}
		total := 0
		for _, l := range latencies {
			total += l
		}
		return total / len(latencies)
	}

	sort.Ints(latencies)
	var avgs [4]float64
	for i, rate := range rates {
		n := int(math.Ceil((1 - rate) * float64(len(latencies))))
		avgs[i] = float64(average(latencies[len(latencies)-n:])) / 1000000
	}
	LogTime(fmt.Sprintf("%.2f\t%.2f\t%.2f\t%.2f\t%s", avgs[0], avgs[1], avgs[2], avgs[3], strings.Join(variables, "\t")))
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

func Itoas(nums ...int) []string {
	r := []string{}
	for _, n := range nums {
		r = append(r, fmt.Sprintf("%d", n))
	}
	return r
}

func Ftoas(nums ...float64) []string {
	r := []string{}
	for _, n := range nums {
		r = append(r, fmt.Sprintf("%0.4f", n))
	}
	return r
}
