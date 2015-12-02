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
	shortTestPeriod = 2 * time.Second
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
	// This is just the plan, we'll calculate the real frequency
	eventFrequency = []int{
		1, // events/second
		2,
		5,
		10,
		20,
		50,
		100,
		200,
	}
	eventRoutines = []int{
		1,
		1,
		2,
		3,
		4,
		10,
		50,
		100,
	}
)

// For docker event loss rate
var (
	defaultEventFrequency = 100
	defaultEventRoutines  = 100 //50
	timesForEachPeriod    = 3
	testPeriodList        = []time.Duration{
		10 * time.Second,
		20 * time.Second,
		30 * time.Second,
		1 * time.Minute,
		2 * time.Minute,
	}
)
