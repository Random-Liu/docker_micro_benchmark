#!/bin/bash

MIN_DOCKER_NUMBER=6000

DOCKER_NUMBER=`docker ps -a | wc -l`
if [ $DOCKER_NUMBER -lt $MIN_DOCKER_NUMBER ]
then
  echo "Benchmark with different container number"
  LC_ALL=C sar -rubwS -P ALL 1 > sar_benchmark_varies_containers.txt &
  SAR_PID=$!
  ./docker_micro_benchmark -c > latency_benchmark_varies_containers.txt
  kill $SAR_PID
fi

echo "Benchmark with different period"
LC_ALL=C sar -rubwS -P ALL 1 > sar_benchmark_varies_period.txt &
SAR_PID=$!
./docker_micro_benchmark -p > latency_benchmark_varies_period.txt
kill $SAR_PID
