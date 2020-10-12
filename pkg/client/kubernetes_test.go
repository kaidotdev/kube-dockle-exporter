package client_test

import (
	"context"
	"errors"
	"fmt"
	"kube-dockle-exporter/pkg/client"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	apiV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type (
	deploymentListFunc  = func(context.Context, metaV1.ListOptions) (*apiV1.DeploymentList, error)
	statefulSetListFunc = func(context.Context, metaV1.ListOptions) (*apiV1.StatefulSetList, error)
	daemonSetListFunc   = func(context.Context, metaV1.ListOptions) (*apiV1.DaemonSetList, error)
)

type kubernetesClientsetMock struct {
	kubernetes.Interface
	fakeDeploymentList            deploymentListFunc
	wantFakeDeploymentListCalled  int
	fakeDeploymentListCalled      int
	fakeStatefulSetList           statefulSetListFunc
	wantFakeStatefulSetListCalled int
	fakeStatefulSetListCalled     int
	fakeDaemonSetList             daemonSetListFunc
	wantFakeDaemonSetListCalled   int
	fakeDaemonSetListCalled       int
}

func (m *kubernetesClientsetMock) assert(t *testing.T) {
	if diff := cmp.Diff(m.wantFakeDeploymentListCalled, m.fakeDeploymentListCalled); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(m.wantFakeStatefulSetListCalled, m.fakeStatefulSetListCalled); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(m.wantFakeDaemonSetListCalled, m.fakeDaemonSetListCalled); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func (m *kubernetesClientsetMock) AppsV1() appsV1.AppsV1Interface {
	return &appsV1Mock{
		fakeDeploymentList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DeploymentList, error) {
			m.fakeDeploymentListCalled++
			return m.fakeDeploymentList(ctx, opts)
		},
		fakeStatefulSetList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.StatefulSetList, error) {
			m.fakeStatefulSetListCalled++
			return m.fakeStatefulSetList(ctx, opts)
		},
		fakeDaemonSetList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DaemonSetList, error) {
			m.fakeDaemonSetListCalled++
			return m.fakeDaemonSetList(ctx, opts)
		},
	}
}

type appsV1Mock struct {
	appsV1.AppsV1Interface
	fakeDeploymentList  deploymentListFunc
	fakeStatefulSetList statefulSetListFunc
	fakeDaemonSetList   daemonSetListFunc
}

func (m *appsV1Mock) Deployments(namespace string) appsV1.DeploymentInterface {
	return &deploymentMock{
		fakeList: m.fakeDeploymentList,
	}
}

type deploymentMock struct {
	appsV1.DeploymentInterface
	fakeList deploymentListFunc
}

func (m *deploymentMock) List(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DeploymentList, error) {
	return m.fakeList(ctx, opts)
}

func (m *appsV1Mock) StatefulSets(namespace string) appsV1.StatefulSetInterface {
	return &statefulSetMock{
		fakeList: m.fakeStatefulSetList,
	}
}

type statefulSetMock struct {
	appsV1.StatefulSetInterface
	fakeList statefulSetListFunc
}

func (m *statefulSetMock) List(ctx context.Context, opts metaV1.ListOptions) (*apiV1.StatefulSetList, error) {
	return m.fakeList(ctx, opts)
}

func (m *appsV1Mock) DaemonSets(namespace string) appsV1.DaemonSetInterface {
	return &daemonSetMock{
		fakeList: m.fakeDaemonSetList,
	}
}

type daemonSetMock struct {
	appsV1.DaemonSetInterface
	fakeList daemonSetListFunc
}

func (m *daemonSetMock) List(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DaemonSetList, error) {
	return m.fakeList(ctx, opts)
}

func TestKubernetesClientContainers(t *testing.T) {
	fakeDeployment := apiV1.Deployment{
		Spec: apiV1.DeploymentSpec{
			Template: coreV1.PodTemplateSpec{
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Image: "deployment",
						},
					},
				},
			},
		},
	}
	fakeStatefulSet := apiV1.StatefulSet{
		Spec: apiV1.StatefulSetSpec{
			Template: coreV1.PodTemplateSpec{
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Image: "statefulSet",
						},
					},
				},
			},
		},
	}
	fakeDaemonSet := apiV1.DaemonSet{
		Spec: apiV1.DaemonSetSpec{
			Template: coreV1.PodTemplateSpec{
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Image: "daemonSet",
						},
					},
				},
			},
		},
	}

	type want struct {
		first []v1.Container
	}

	tests := []struct {
		name            string
		receiver        *client.KubernetesClient
		want            want
		wantErrorString string
		optsFunction    func(interface{}) cmp.Option
	}{
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			&client.KubernetesClient{
				Inner: &kubernetesClientsetMock{
					fakeDeploymentList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DeploymentList, error) {
						return &apiV1.DeploymentList{
							Items: []apiV1.Deployment{
								fakeDeployment,
							},
						}, nil
					},
					wantFakeDeploymentListCalled: 1,
					fakeStatefulSetList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.StatefulSetList, error) {
						return &apiV1.StatefulSetList{
							Items: []apiV1.StatefulSet{
								fakeStatefulSet,
							},
						}, nil
					},
					wantFakeStatefulSetListCalled: 1,
					fakeDaemonSetList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DaemonSetList, error) {
						return &apiV1.DaemonSetList{
							Items: []apiV1.DaemonSet{
								fakeDaemonSet,
							},
						}, nil
					},
					wantFakeDaemonSetListCalled: 1,
				},
			},
			want{
				[]v1.Container{
					{
						Image: "deployment",
					},
					{
						Image: "statefulSet",
					},
					{
						Image: "daemonSet",
					},
				},
			},
			"",
			func(got interface{}) cmp.Option {
				return nil
			},
		},
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			&client.KubernetesClient{
				Inner: &kubernetesClientsetMock{
					fakeDeploymentList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DeploymentList, error) {
						return nil, errors.New("fake")
					},
					wantFakeDeploymentListCalled:  1,
					wantFakeStatefulSetListCalled: 0,
					wantFakeDaemonSetListCalled:   0,
				},
			},
			want{
				nil,
			},
			"could not get deployment: fake",
			func(got interface{}) cmp.Option {
				return nil
			},
		},
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			&client.KubernetesClient{
				Inner: &kubernetesClientsetMock{
					fakeDeploymentList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DeploymentList, error) {
						return &apiV1.DeploymentList{
							Items: []apiV1.Deployment{
								fakeDeployment,
							},
						}, nil
					},
					wantFakeDeploymentListCalled: 1,
					fakeStatefulSetList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.StatefulSetList, error) {
						return nil, errors.New("fake")
					},
					wantFakeStatefulSetListCalled: 1,
					wantFakeDaemonSetListCalled:   0,
				},
			},
			want{
				nil,
			},
			"could not get stateful set: fake",
			func(got interface{}) cmp.Option {
				return nil
			},
		},
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			&client.KubernetesClient{
				Inner: &kubernetesClientsetMock{
					fakeDeploymentList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DeploymentList, error) {
						return &apiV1.DeploymentList{
							Items: []apiV1.Deployment{
								fakeDeployment,
							},
						}, nil
					},
					wantFakeDeploymentListCalled: 1,
					fakeStatefulSetList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.StatefulSetList, error) {
						return &apiV1.StatefulSetList{
							Items: []apiV1.StatefulSet{
								fakeStatefulSet,
							},
						}, nil
					},
					wantFakeStatefulSetListCalled: 1,
					fakeDaemonSetList: func(ctx context.Context, opts metaV1.ListOptions) (*apiV1.DaemonSetList, error) {
						return nil, errors.New("fake")
					},
					wantFakeDaemonSetListCalled: 1,
				},
			},
			want{
				nil,
			},
			"could not get daemon set: fake",
			func(got interface{}) cmp.Option {
				return nil
			},
		},
	}
	for _, tt := range tests {
		name := tt.name
		receiver := tt.receiver
		want := tt.want
		wantErrorString := tt.wantErrorString
		optsFunction := tt.optsFunction
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := receiver.Containers()
			receiver.Inner.(*kubernetesClientsetMock).assert(t)
			if diff := cmp.Diff(want.first, got, optsFunction(got)); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}

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
