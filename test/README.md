# Test Environment

## Run nginx
```bash
docker-compose up -d
```

## Run snitch
```bash
docker run --rm -v "${PWD}/log:/var/log/:ro" -v "${PWD}/config.yml:/etc/snitch/config.yml:ro" -e "HOST=localhost" xjewer/snitch:v0.7 -config /etc/snitch/config.yml -statsd 127.0.0.1:8125 -buffer 10
```

## Test
```bash
curl http://127.0.0.1:9130
curl http://127.0.0.1:9130/1
```
