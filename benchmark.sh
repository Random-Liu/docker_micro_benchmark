#!/bin/bash
RESULT=result

if [ ! -d $RESULT ]
then
	mkdir $RESULT
fi

DOCKER_NUMBER=`docker ps -a | wc -l`
DOCKER_NUMBER=`expr $DOCKER_NUMBER - 1`
if [ $DOCKER_NUMBER -eq 0 ]
then
	echo "Benchmark with different container numbers"
	LC_ALL=C sar -rubwS -P ALL 1 > $RESULT/sar_benchmark_varies_containers.txt &
	SAR_PID=$!
	./docker_micro_benchmark -c > $RESULT/latency_benchmark_varies_containers.txt
	kill $SAR_PID
else
	echo "Exsiting Docker number: $DOCKER_NUMBER, skip benchmark with different container numbers"
fi

echo "Benchmark with different periods"
LC_ALL=C sar -rubwS -P ALL 1 > $RESULT/sar_benchmark_varies_period.txt &
SAR_PID=$!
./docker_micro_benchmark -p > $RESULT/latency_benchmark_varies_period.txt
kill $SAR_PID

echo "Benchmark with different routine numbers"
LC_ALL=C sar -rubwS -P ALL 1 > $RESULT/sar_benchmark_varies_routines.txt &
SAR_PID=$!
./docker_micro_benchmark -r > $RESULT/latency_benchmark_varies_routines.txt
kill $SAR_PID

echo "Benchmark event stream"
LC_ALL=C sar -rubwS -P ALL 1 > $RESULT/sar_benchmark_event_stream.txt &
SAR_PID=$!
shell/pidstat-grapher.py -a docker,docker_micro_benchmark -d $RESULT > /dev/null &
PIDSTAT_PID=$!
./docker_micro_benchmark -e > $RESULT/latency_benchmark_event_stream.txt 
kill $SAR_PID
kill $PIDSTAT_PID
