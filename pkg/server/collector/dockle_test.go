package collector_test

import (
	"context"
	"errors"
	"fmt"
	"kube-dockle-exporter/pkg/server/collector"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
)

func getRecursiveStructReflectValue(rv reflect.Value) []reflect.Value {
	var values []reflect.Value
	switch rv.Kind() {
	case reflect.Ptr:
		values = append(values, getRecursiveStructReflectValue(rv.Elem())...)
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			values = append(values, getRecursiveStructReflectValue(rv.Index(i))...)
		}
	case reflect.Map:
		for _, k := range rv.MapKeys() {
			values = append(values, getRecursiveStructReflectValue(rv.MapIndex(k))...)
		}
	case reflect.Struct:
		values = append(values, reflect.New(rv.Type()).Elem())
		for i := 0; i < rv.NumField(); i++ {
			values = append(values, getRecursiveStructReflectValue(rv.Field(i))...)
		}
	default:
	}
	return values
}

func TestDockleCollectorDescribe(t *testing.T) {
	tests := []struct {
		name         string
		receiver     *collector.DockleCollector
		in           chan *prometheus.Desc
		want         *prometheus.Desc
		optsFunction func(interface{}) cmp.Option
	}{
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			collector.NewDockleCollector(
				&loggerMock{
					wantFakeErrorfCalled: 0,
					wantFakeInfofCalled:  0,
					wantFakeDebugfCalled: 0,
				},
				&kubernetesClientMock{
					wantFakeContainersCalled: 0,
				},
				&dockleClientMock{
					wantFakeDoCalled: 0,
				},
				1,
			),
			make(chan *prometheus.Desc, 1),
			prometheus.NewDesc(
				"dockle_cis_benchmarks_total",
				"CIS benchmarks executed by dockle",
				[]string{"image", "code", "level"},
				nil,
			),
			func(got interface{}) cmp.Option {
				switch v := got.(type) {
				case *prometheus.Desc:
					return cmp.AllowUnexported(*v)
				default:
					return nil
				}
			},
		},
	}
	for _, tt := range tests {
		name := tt.name
		receiver := tt.receiver
		in := tt.in
		want := tt.want
		optsFunction := tt.optsFunction
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			receiver.Describe(in)
			got := <-in
			receiver.Logger.(*loggerMock).assert(t)
			receiver.KubernetesClient.(*kubernetesClientMock).assert(t)
			receiver.DockleClient.(*dockleClientMock).assert(t)
			if diff := cmp.Diff(want, got, optsFunction(got)); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestDockleCollectorCollect(t *testing.T) {
	tests := []struct {
		name         string
		receiver     *collector.DockleCollector
		in           chan prometheus.Metric
		want         prometheus.Metric
		optsFunction func(interface{}) cmp.Option
	}{
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			collector.NewDockleCollector(
				&loggerMock{
					wantFakeErrorfCalled: 0,
					wantFakeInfofCalled:  0,
					wantFakeDebugfCalled: 0,
				},
				&kubernetesClientMock{
					fakeContainers: func() ([]v1.Container, error) {
						return []v1.Container{
							{
								Image: "fake",
							},
						}, nil
					},
					wantFakeContainersCalled: 1,
				},
				&dockleClientMock{
					fakeDo: func(ctx context.Context, image string) ([]byte, error) {
						return []byte(`{"Target":"fake","Details":[{"code":"fake"}]}`), nil
					},
					wantFakeDoCalled: 1,
				},
				1,
			),
			make(chan prometheus.Metric, 1),
			func() prometheus.Gauge {
				gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "dockle",
					Name:      "cis_benchmarks_total",
					Help:      "CIS benchmarks executed by dockle",
				}, []string{"image", "code", "level"})
				labels := []string{
					"fake",
					"fake",
					"",
				}
				gaugeVec.WithLabelValues(labels...).Set(1)
				gauge, err := gaugeVec.GetMetricWithLabelValues(labels...)
				if err != nil {
					t.Fatal()
				}
				return gauge
			}(),
			func(got interface{}) cmp.Option {
				switch got.(type) {
				case prometheus.Metric:
					deref := func(v interface{}) interface{} {
						return reflect.ValueOf(v).Elem().Interface()
					}
					v := deref(got)
					switch reflect.TypeOf(v).Name() {
					case "gauge":
						var opts cmp.Options
						for _, rv := range getRecursiveStructReflectValue(reflect.ValueOf(v)) {
							switch rv.Type().Name() {
							case "selfCollector":
								opts = append(opts, cmpopts.IgnoreUnexported(rv.Interface()))
							default:
								opts = append(opts, cmp.AllowUnexported(rv.Interface()))
							}
						}
						return opts
					default:
						return nil
					}
				default:
					return nil
				}
			},
		},
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			collector.NewDockleCollector(
				&loggerMock{
					fakeErrorf: func(format string, v ...interface{}) {
						want := "Failed to scan: failed to get containers: fake\n"
						got := fmt.Sprintf(format, v...)
						if diff := cmp.Diff(want, got); diff != "" {
							t.Errorf("(-want +got):\n%s", diff)
						}
					},
					wantFakeErrorfCalled: 1,
					wantFakeInfofCalled:  0,
					wantFakeDebugfCalled: 0,
				},
				&kubernetesClientMock{
					fakeContainers: func() ([]v1.Container, error) {
						return nil, errors.New("fake")
					},
					wantFakeContainersCalled: 1,
				},
				&dockleClientMock{
					wantFakeDoCalled: 0,
				},
				1,
			),
			func() chan prometheus.Metric {
				ch := make(chan prometheus.Metric, 1)
				close(ch)
				return ch
			}(),
			nil,
			func(got interface{}) cmp.Option {
				return nil
			},
		},
	}
	for _, tt := range tests {
		name := tt.name
		receiver := tt.receiver
		in := tt.in
		want := tt.want
		optsFunction := tt.optsFunction
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			receiver.StartLoop(ctx, 10*time.Millisecond)
			time.Sleep(15 * time.Millisecond)
			receiver.Collect(in)
			got := <-in
			cancel()
			receiver.Logger.(*loggerMock).assert(t)
			receiver.KubernetesClient.(*kubernetesClientMock).assert(t)
			receiver.DockleClient.(*dockleClientMock).assert(t)
			if diff := cmp.Diff(want, got, optsFunction(got)); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestDockleCollectorScan(t *testing.T) {
	type in struct {
		first context.Context
	}

	tests := []struct {
		name            string
		receiver        *collector.DockleCollector
		in              in
		wantErrorString string
	}{
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			collector.NewDockleCollector(
				&loggerMock{
					wantFakeErrorfCalled: 0,
					wantFakeInfofCalled:  0,
					wantFakeDebugfCalled: 0,
				},
				&kubernetesClientMock{
					fakeContainers: func() ([]v1.Container, error) {
						return nil, errors.New("fake")
					},
					wantFakeContainersCalled: 1,
				},
				&dockleClientMock{
					wantFakeDoCalled: 0,
				},
				1,
			),
			in{
				context.Background(),
			},
			"failed to get containers: fake",
		},
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			collector.NewDockleCollector(
				&loggerMock{
					fakeErrorf: func(format string, v ...interface{}) {
						want := "Failed to detect CIS benchmark at fake: fake\n"
						got := fmt.Sprintf(format, v...)
						if diff := cmp.Diff(want, got); diff != "" {
							t.Errorf("(-want +got):\n%s", diff)
						}
					},
					wantFakeErrorfCalled: 1,
					wantFakeInfofCalled:  0,
					wantFakeDebugfCalled: 0,
				},
				&kubernetesClientMock{
					fakeContainers: func() ([]v1.Container, error) {
						return []v1.Container{
							{
								Image: "fake",
							},
						}, nil
					},
					wantFakeContainersCalled: 1,
				},
				&dockleClientMock{
					fakeDo: func(ctx context.Context, image string) ([]byte, error) {
						return nil, errors.New("fake")
					},
					wantFakeDoCalled: 1,
				},
				1,
			),
			in{
				context.Background(),
			},
			"",
		},
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			collector.NewDockleCollector(
				&loggerMock{
					fakeErrorf: func(format string, v ...interface{}) {
						want := "Failed to parse dockle response at fake: invalid character 'k' in literal false (expecting 'l')\n"
						got := fmt.Sprintf(format, v...)
						if diff := cmp.Diff(want, got); diff != "" {
							t.Errorf("(-want +got):\n%s", diff)
						}
					},
					wantFakeErrorfCalled: 1,
					wantFakeInfofCalled:  0,
					wantFakeDebugfCalled: 0,
				},
				&kubernetesClientMock{
					fakeContainers: func() ([]v1.Container, error) {
						return []v1.Container{
							{
								Image: "fake",
							},
						}, nil
					},
					wantFakeContainersCalled: 1,
				},
				&dockleClientMock{
					fakeDo: func(ctx context.Context, image string) ([]byte, error) {
						return []byte("fake"), nil
					},
					wantFakeDoCalled: 1,
				},
				1,
			),
			in{
				context.Background(),
			},
			"",
		},
	}
	for _, tt := range tests {
		name := tt.name
		receiver := tt.receiver
		in := tt.in
		wantErrorString := tt.wantErrorString
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := receiver.Scan(in.first)
			receiver.Logger.(*loggerMock).assert(t)
			receiver.KubernetesClient.(*kubernetesClientMock).assert(t)
			receiver.DockleClient.(*dockleClientMock).assert(t)

			if err == nil {
				if diff := cmp.Diff(wantErrorString, ""); diff != "" {
					t.Errorf("(-want +got):\n%s", diff)
				}
			} else {
				gotErrorString := err.Error()
				if diff := cmp.Diff(wantErrorString, gotErrorString); diff != "" {
					t.Errorf("(-want +got):\n%s", diff)
				}
			}
		})
	}
}
