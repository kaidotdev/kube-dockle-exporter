# KubeDockleExporter

KubeDockleExporter is Prometheus Exporter that collects CIS benchmarks executed by [goodwithtech/dockle](https://github.com/goodwithtech/dockle) in the kubernetes cluster.

## Installation

```shell
$ kubectl apply -k manifests
```

## Usage

```shell
$ curl http://kube-dockle-exporter:9090/metrics | grep dockle_cis_benchmarks_total | head -n 10
# HELP dockle_cis_benchmarks_total CIS benchmarks executed by dockle
# TYPE dockle_cis_benchmarks_total gauge
dockle_cis_benchmarks_total{code="CIS-DI-0001",image="alpine/socat:1.0.5",level="WARN"} 1
dockle_cis_benchmarks_total{code="CIS-DI-0001",image="docker.elastic.co/elasticsearch/elasticsearch-oss:7.9.2",level="WARN"} 1
dockle_cis_benchmarks_total{code="CIS-DI-0001",image="docker.io/cilium/cilium:v1.8.4",level="WARN"} 1
dockle_cis_benchmarks_total{code="CIS-DI-0001",image="docker.io/cilium/operator-generic:v1.8.4",level="WARN"} 1
dockle_cis_benchmarks_total{code="CIS-DI-0001",image="docker.io/falcosecurity/falco:0.25.0",level="WARN"} 1
dockle_cis_benchmarks_total{code="CIS-DI-0001",image="docker.io/istio/proxyv2:1.6.8",level="WARN"} 1
dockle_cis_benchmarks_total{code="CIS-DI-0001",image="docker.io/jaegertracing/all-in-one:1.16",level="WARN"} 1
dockle_cis_benchmarks_total{code="CIS-DI-0001",image="docker.io/kennethreitz/httpbin",level="WARN"} 1
```

## How to develop

### `skaffold dev`

```sh
$ make dev
```

### Test

```sh
$ make test
```

### Lint

```sh
$ make lint
```
