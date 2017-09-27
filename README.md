## Snitch

This tool allows to parse log files and send statistics to statsd endpoint

### Configuration

Snith has a particular configuration structure in yml:

```
sources:
- source: name of source
  file: ./test/log/test2.log
  noFollow: false
  mustExists: false
  reOpen: true
  prefix: "%HOST%"
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
* {string} `prefix` - statsd key prefix
* {string} `delimiter` - column delimiter in log files, it is reasonable to use `\t` delimiter in log files

* {[]key} `keys` - list of keys, which would be send to statsd 
    * {string} `key` - statsd metric key, $N - means column position
    * {boolean} `count` - boolean, means  
    * {string} `timing` - column with time (time should be in seconds, like 0.001 - means 1ms)  
    * {delimiter} `delimiter` - delimiter into the column, 
            e.g. nginx can to write in one column a few values from upstream the request was in   


Snitch allows handle a few sources per instance

### Docker build

```
docker run --rm -v "$(pwd):/src" -v /var/run/docker.sock:/var/run/docker.sock -e MAIN_PATH=cmd/snitch xjewer/golang-builder:v1.9 xjewer/snitch:{tag}
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
