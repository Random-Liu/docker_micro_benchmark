package main

import (
	"time"
)

// TODO Use configuration file later
var (
	periods = []time.Duration{
		0 * time.Second,
		200 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
		2 * time.Second,
		5 * time.Second,
		10 * time.Second,
	}
	defaultPeriod   = 200 * time.Millisecond
	shortTestPeriod = 10 * time.Second
	longTestPeriod  = 50 * time.Second
	routines        = []int{10, 50, 100}
)

// For varies container number benchmark
var (
	deadContainers = []int{
		100,
		500,
		1000,
		2000,
		5000,
	}
	aliveContainers = []int{
		50,
		100,
		500,
		1000,
	}
)
