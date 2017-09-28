## Snitch

[![Travis](https://img.shields.io/travis/xjewer/snitch.svg)](https://travis-ci.org/xjewer/snitch)
[![Coveralls](https://img.shields.io/coveralls/xjewer/snitch.svg)](https://coveralls.io/github/xjewer/snitch)
[![Go Report Card](https://goreportcard.com/badge/github.com/xjewer/snitch)](https://goreportcard.com/report/github.com/xjewer/snitch)
[![GitHub release](https://img.shields.io/github/release/xjewer/snitch.svg)](https://github.com/xjewer/snitch/releases)
[![Docker Automated build](https://img.shields.io/docker/automated/xjewer/snitch.svg)](https://hub.docker.com/r/xjewer/snitch/)
[![ImageLayers Layers](https://img.shields.io/imagelayers/layers/xjewer/snitch/latest.svg)](https://hub.docker.com/r/xjewer/snitch/)
[![ImageLayers Size](https://img.shields.io/imagelayers/image-size/xjewer/snitch/latest.svg)](https://hub.docker.com/r/xjewer/snitch/)

This tool allows to parse log files and send statistics to statsd endpoint

### Startup options

* {string} `config` - configuration file
* {string} `statsd` - statsd endpoint
* {string} `prefix` - statsd global key prefix, e.g. `balancer.` 
* {int} `buffer` - buffer interval for metrics, `0` default: one metric - one request 

### Configuration file

Snith has a particular configuration structure in yml:

```
sources:
- source: name of source
  file: ./test/log/test2.log
  noFollow: false
  mustExists: false
  reOpen: true
  prefix: "balancer.%HOST%"
  delimiter: "\t"
  keys:
    - key: All.$3.$6
      count: true
      timing: $4
      delimiter: " : "
    - key: All.$3.
      timing: $10
      delimiter: " : "
```

where is:

* {string} `source`  - source name
* {string} `file` - where the log file is
* {boolean} `noFollow` - means no follow new lines that are written in the log file, if true
* {boolean} `mustExists` - log file have to be existed, if true
* {boolean} `reOpen` - re-open file, if true (e.g. log rotation)
* {string} `prefix` - statsd key prefix, `%HOST%` is used as substitution a hostname
* {string} `delimiter` - column delimiter in log files, it is reasonable to use `\t` delimiter in log files

* {[]key} `keys` - list of keys, which would be send to statsd 
    * {string} `key` - statsd metric key, $N - means column position
    * {boolean} `count` - boolean, means  
    * {string} `timing` - column with time (time should be in seconds, like 0.001 - means 1ms)  
    * {string} `delimiter` - delimiter into the column, 
            e.g. nginx can to write in one column a few values from upstream the request was in   


Snitch allows handle a few sources per instance

### Docker build

Docker build uses [multistage build](https://docs.docker.com/engine/userguide/eng-image/multistage-build/)

```
docker build -t xjewer/snitch:{tag} .
```


### Run


```
docker run --rm -v "/var/log/nginx/balancer/:/var/log/:ro" xjewer/snitch:0.1 -file /var/log/for_script.access.log -statsd graphite.local:8125 -buffer 10
```

### Benchmarks

Benchmark line parsers is:

```
goos: darwin
goarch: amd64
pkg: github.com/xjewer/snitch
Benchmark_parser-4   	 2000000	       946 ns/op	     416 B/op	       8 allocs/op
PASS
ok  	github.com/xjewer/snitch	2.890s

```
