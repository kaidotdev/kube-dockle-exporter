package collector_test

import (
	"context"
	"kube-dockle-exporter/pkg/server/collector"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
)

type loggerMock struct {
	collector.ILogger
	fakeErrorf           func(format string, v ...interface{})
	wantFakeErrorfCalled int
	fakeErrorfCalled     int
	fakeInfof            func(format string, v ...interface{})
	wantFakeInfofCalled  int
	fakeInfofCalled      int
	fakeDebugf           func(format string, v ...interface{})
	wantFakeDebugfCalled int
	fakeDebugfCalled     int
}

func (m *loggerMock) assert(t *testing.T) {
	if diff := cmp.Diff(m.wantFakeErrorfCalled, m.fakeErrorfCalled); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(m.wantFakeInfofCalled, m.fakeInfofCalled); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(m.wantFakeDebugfCalled, m.fakeDebugfCalled); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func (m *loggerMock) Errorf(format string, v ...interface{}) {
	m.fakeErrorfCalled++
	m.fakeErrorf(format, v...)
}

func (m *loggerMock) Infof(format string, v ...interface{}) {
	m.fakeInfofCalled++
	m.fakeInfof(format, v...)
}

func (m *loggerMock) Debugf(format string, v ...interface{}) {
	m.fakeDebugfCalled++
	m.fakeDebugf(format, v...)
}

type kubernetesClientMock struct {
	collector.IKubernetesClient
	fakeContainers           func() ([]v1.Container, error)
	wantFakeContainersCalled int
	fakeContainersCalled     int
}

func (m *kubernetesClientMock) assert(t *testing.T) {
	if diff := cmp.Diff(m.wantFakeContainersCalled, m.fakeContainersCalled); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func (m *kubernetesClientMock) Containers() ([]v1.Container, error) {
	m.fakeContainersCalled++
	return m.fakeContainers()
}

type dockleClientMock struct {
	collector.IDockleClient
	fakeDo           func(context.Context, string) ([]byte, error)
	wantFakeDoCalled int
	fakeDoCalled     int
}

func (m *dockleClientMock) assert(t *testing.T) {
	if diff := cmp.Diff(m.wantFakeDoCalled, m.fakeDoCalled); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func (m *dockleClientMock) Do(ctx context.Context, image string) ([]byte, error) {
	m.fakeDoCalled++
	return m.fakeDo(ctx, image)
}
