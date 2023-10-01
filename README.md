# k6-go

k6-go makes it easy to run load tests written in golang with k6. To achieve this it uses [xk6](https://github.com/grafana/xk6) and the [xk6-g0](https://github.com/szkiba/xk6-g0) extension.

## build

```bash
go build cmd/k6-go/k6-go.go
```

## example

[petstore-load-k6](https://github.com/bandorko/petstore-load-k6) is an example load test project written in golang

**asciicast (example run)**
[![asciicast](https://asciinema.org/a/alS0Z1IRUJzp71gIzbuC1djQi.svg)](https://asciinema.org/a/alS0Z1IRUJzp71gIzbuC1djQi)