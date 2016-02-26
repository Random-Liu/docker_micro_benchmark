Docker Micro Benchmark
======================
Docker micro benchmark is a tool aimed at benchmarking docker operations
which are critical to K8s performance, such as `docker ps`, `docker inspect` etc.

Description
------------
Docker micro benchmark benchmarks the following docker operations:
  * **`docker ps [-a]`**: K8s does periodically `docker ps -a` to detect container
    state changes, so its performance is crutial to K8s.
  * **`docker inspect`**: K8s does `docker inspect` to get detailed information of
    specific container when it finds out a container state is changed. So inspect 
    is also relatively frequent.
  * **`docker create` & `docker start`**: The performance of `docker create` and
    `docker start` are important for K8s when creating pods, especially for batch
    creation.
  * **`docker stop` & `docker remove`**: The same with above.

Docker micro benchmark supports 4 kinds of benchmarks:
  * Benchmark `docker ps`, `docker ps -a` and `docker inspect` with different number
    of dead and alive containers.
  * Benchmark `docker ps`, `docker ps -a` and `docker inspect` with different operation
    interval.
  * Benchmark `docker ps -a` and `docker inspect` with different number of goroutines.
  * Benchmark `docker create & docker start` and `docker stop & docker remove` with
    different operation rate limit.

Instructions
------------
#### Dependency
* golang
* sysstat
* gnuplot

#### Build
`go build github.com/random-liu/docker_micro_benchmark`

#### Usage
`benchmark.sh` is the script starting the benchmark.

`Usage : benchmark.sh -[o|c|i|r]`:
  * **-o**: Run `docker create/start/stop/remove` benchmark.
  * **-c**: Run `docker list/inspect` with different number of containers benchmark.
    *Notice that containers created in this benchmark won't be removed before
    benchmark is over.* That's intendend because creating containers is really
    slow, it's better to reuse these containers in the following benchmark. If
    you want to remove all the containers, just run `shell/kill_all_container.sh`.
  * **-i**: Run `docker list/inspect` with different operation interval.
  * **-r**: Run `docker list/inspect` with different number of goroutines.

You can run `benchmark.sh` with multiple options, the script will run benchmarks
corresponding to each of the options one by one. *Notice that it's better to put
`-c` in front of `-i` and `-r`, so that they can reuse the containers created in
`-c` benchmark.*

#### Result
After the benchmark finishes, all benchmark result will be generated in `result/`
directory. Different benchmark results locate in different sub-directory.

There are two forms of benchmark result:
* **\*.dat**: Table formatted text result.
* **\*.png**: Graph result auto-generated from the text result. Notice that there
    is no graph for `sar` result. If you want to analyse the data in `sar` result,
    you can use another tool [kSar](https://sourceforge.net/projects/ksar/).

Result
------
Here are links of some previous benchmark results:
* https://github.com/kubernetes/kubernetes/issues/16110#issuecomment-180510177
* https://docs.google.com/document/d/1d5xaYW3oVnzTgjlcRSnj5aJQusRkrQotXF2zBoTnNkw/edit


