# k6-go

k6-go makes it easy to run load tests written in golang with k6. To achieve this it uses [xk6](https://github.com/grafana/xk6) and the [xk6-g0](https://github.com/szkiba/xk6-g0) extension.

## the problem

Although k6's officially supported scripting language is JavaScript, since xk6-g0 was born, writing scripts in golang became possible. xk6-go can handle simple go scripts out of the box with the help of [yaegi](https://github.com/traefik/yaegi) interpreter, but if you want to use subpackages, or 3rd party packages (like a generated openapi client), the situation became harder. It is possible, to install additional packages to the yaegi interpreter (described [here](https://github.com/szkiba/xk6-g0#extending-xk6-g0)), but the procedure is not developer friendly.

## the solution

k6-go detects the dependencies of the go script, and generates the extension, that uses the RegisterExports function of xk6-go, to build the needed packages into a custom k6 binary on the fly.

## installing k6-go

```bash
go install github.com/bandorko/k6-go@latest
```

## running test

k6-go runs the custom k6 binary after building it, so you can use k6-go with the same parameters as k6.
```bash
k6-go run -i 10 -u 3 script.go
k6-go help run
```


If you want to build only the custom k6 binary, then you can use the build subcommand
```bash
k6-go build --output custom-k6 script.go
```

## example

[petstore-load-k6](https://github.com/bandorko/petstore-load-k6) is an example load test project written in golang

**asciicast (example run)**
[![asciicast](https://asciinema.org/a/alS0Z1IRUJzp71gIzbuC1djQi.svg)](https://asciinema.org/a/alS0Z1IRUJzp71gIzbuC1djQi)