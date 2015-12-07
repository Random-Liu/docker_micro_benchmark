#!/bin/bash
RESULT=result
DOCKER_MICRO_BENCHMARK=./docker_micro_benchmark
STAT_TOOL=pidstat

usage () {
  echo 'Usage : $0 -[c|p|r|e|l]'
  exit
}

# $1 parameter, $2 file suffix
doBenchmark() {
  if [ ! -d $2 ]; then
    mkdir $RESULT/$2
  fi
  LC_ALL=C sar -rubwS -P ALL 1 > $RESULT/$2/sar_benchmark_$2.txt &
  SAR_PID=$!
  $DOCKER_MICRO_BENCHMARK $1 > $RESULT/$2/result_benchmark_$2.txt &
  BENCHMARK_PID=$!
  $STAT_TOOL -p $BENCHMARK_PID 1 > $RESULT/$2/cpu_benchmark_$2.txt &
  DOCKER_PID=`ps -ef | awk '$8=="/usr/bin/docker" {print $2}'`
  $STAT_TOOL -p $DOCKER_PID 1 > $RESULT/$2/cpu_docker_daemon_$2.txt &
  DOCKER_PIDSTAT=$!
  wait $BENCHMARK_PID
  kill $SAR_PID
  kill $DOCKER_PIDSTAT
  kill $SAR_PID
}

if [ ! -d $RESULT ]; then
  mkdir $RESULT
fi

while [ "$1" != "" ]; do
  case $1 in
    -c )
      DOCKER_NUMBER=`docker ps -a | wc -l`
      DOCKER_NUMBER=`expr $DOCKER_NUMBER - 1`
      if [ $DOCKER_NUMBER -ne 0 ]; then
        shell/kill_all_dockers.sh
      fi 
      echo "Benchmark with different container numbers"
      doBenchmark $1 varies_containers
      shift
      ;;
    -p )
      echo "Benchmark with different periods"
      doBenchmark $1 varies_period 
      shift
      ;;
    -r )
      echo "Benchmark with different routine numbers"
      doBenchmark $1 varies_routines
      shift
      ;;
    -e )
      echo "Benchmark event stream"
      doBenchmark $1 event_stream
      shift
      ;;
    -l )
      echo "Benchmark event loss rate"
      doBenchmark $1 event_loss_rate
      shift
      ;;
    * )
      usage
      exit 1
      ;;
  esac
done
