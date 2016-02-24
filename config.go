package main

import (
	"time"
)

// Docker configuration
var (
	endpoint = "unix:///var/run/docker.sock"
)

// Period configuration
var (
	defaultPeriod   = 200 * time.Millisecond
	shortTestPeriod = 10 * time.Second
	longTestPeriod  = 50 * time.Second
)

var (
	listPeriods = []time.Duration{
		0 * time.Second,
		50 * time.Millisecond,
		100 * time.Millisecond,
		200 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
		2 * time.Second,
		5 * time.Second,
	}
	inspectPeriods = []time.Duration{
		0 * time.Second,
		1 * time.Millisecond,   // 1000 inspect/second = 100 pods * 10 containers
		2 * time.Millisecond,   // 500 inspect/second = 100 pods * 5 containers = 50 pods * 10 containers
		5 * time.Millisecond,   // 200 inspect/second = 100 pods * 2 containers = 20 pods * 10 containers
		10 * time.Millisecond,  // 100 inspect/second = 100 pods * 1 containers = 10 pods * 10 containers
		50 * time.Millisecond,  // 20 inspect/second = 20 pods * 1 containers = 10 pods * 2 containers
		100 * time.Millisecond, // 10 insepct/second = 10 pods * 1 containers = 5 pods * 2 containers
	}
)

// For container start benchmark
var containerStartConfig = map[string]interface{}{
	"qps": []float64{
		1.0,
		2.0,
		4.0,
		8.0,
		16.0,
		32.0,
		64.0,
	},
	"routine": 100,
}

// For varies container number benchmark
var (
	// aliveContainers * 3
	deadContainers = []int{
		60,
		120,
		300,
		600,
	}
	aliveContainers = []int{
		20,  // 10 * 2
		40,  // 20 * 2
		100, // 50 * 2
		200, // 100 * 2
	}
)

// For varies routine number benchmark
var (
	resyncPeriod         = 1 * time.Second
	routineInspectPeriod = resyncPeriod / 2 // 2 containers/pod
	routines             = []int{
		1,
		5,
		10,
		20,
		50,
		100,
	}
)

// For docker event stream benchmark
var (
	eventInterval = 100 * time.Millisecond
	eventRoutines = []int{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
	}
)

// For docker event loss rate
var (
	defaultEventRoutines = 6
	eventLossTestPeriod  = 1 * time.Hour
)
