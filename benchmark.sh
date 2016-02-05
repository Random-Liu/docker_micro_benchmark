#!/bin/bash
RESULT=result
DOCKER_MICRO_BENCHMARK=./docker_micro_benchmark
STAT_TOOL=pidstat
GNUPLOT=gnuplot
PLOTDIR=plot

usage () {
  echo 'Usage : $0 -[c|p|r|e|l]'
  exit
}

# $1 parameter, $2 directoy name
doBenchmark() {
  RDIR=$RESULT/$2
  if [ ! -d  $RDIR ]; then
    mkdir $RDIR
  fi
  LC_ALL=C sar -rubwS -P ALL 1 > $RDIR/sar_benchmark.dat &
  SAR_PID=$!
  $DOCKER_MICRO_BENCHMARK $1 > $RDIR/result_benchmark.dat &
  BENCHMARK_PID=$!
  $STAT_TOOL -p $BENCHMARK_PID 1 > $RDIR/cpu_benchmark.dat &
  DOCKER_PID=`ps -ef | awk '$8=="/usr/bin/docker" {print $2}'`
  $STAT_TOOL -p $DOCKER_PID 1 > $RDIR/cpu_docker_daemon.dat &
  DOCKER_PIDSTAT=$!
  wait $BENCHMARK_PID
  kill $SAR_PID
  kill $DOCKER_PIDSTAT
  kill $SAR_PID
  cd $RDIR 
  $GNUPLOT ../../$PLOTDIR/cpu_plot
  $GNUPLOT ../../$PLOTDIR/latency_plot
  $GNUPLOT ../../$PLOTDIR/$2/result_plot
  cd - > /dev/null
}

if [ ! -d $RESULT ]; then
  mkdir $RESULT
fi

while [ "$1" != "" ]; do
  case $1 in
    -c )
      CONTAINER_NUMBER=`docker ps -a | wc -l`
      CONTAINER_NUMBER=`expr $CONTAINER_NUMBER - 1`
      if [ $CONTAINER_NUMBER -ne 0 ]; then
        shell/kill_all_dockers.sh > /dev/null
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
# Enable event stream benchmark when gnuplot for event stream is added
#    -e )
#      echo "Benchmark event stream"
#      doBenchmark $1 event_stream
#      shift
#      ;;
#    -l )
#      echo "Benchmark event loss rate"
#      doBenchmark $1 event_loss_rate
#      shift
#      ;;
    * )
      usage
      exit 1
      ;;
  esac
done
