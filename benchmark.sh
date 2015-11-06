#!/bin/bash

LC_ALL=C sar -rubwS -P ALL 1 > test.txt &
SAR_PID=$!
./docker_micro_benchmark
#kill $SAR_PID 
