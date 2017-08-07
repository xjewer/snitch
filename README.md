## Snitch

This tool allows parse log files and sends statistics to statsd endpoint

### Docker build

```
docker run --rm -v "$(pwd):/src" -v /var/run/docker.sock:/var/run/docker.sock -e MAIN_PATH=cmd/snitch xjewer/golang-builder:v1.2 xjewer/snitch:{tag}
```


### Run


```
docker run --rm -v "/var/log/nginx/balancer/:/var/log/:ro" xjewer/snitch:0.1 -file /var/log/for_script.access.log -statsd graphite.local:8125 -buffer 10
```
