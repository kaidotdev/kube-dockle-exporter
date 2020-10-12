module kube-dockle-exporter

go 1.15

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/golangci/golangci-lint v1.26.0 // indirect
	github.com/google/go-cmp v0.4.0
	github.com/gorilla/mux v1.7.4
	github.com/instrumenta/conftest v0.18.2 // indirect
	github.com/instrumenta/kubeval v0.0.0-20190901100547-eae975a0031c // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/prometheus/client_golang v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1 // indirect
	go.opencensus.io v0.22.3
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f
	golang.org/x/sys v0.0.0-20200509044756-6aff5f38e54f
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66 // indirect
)
