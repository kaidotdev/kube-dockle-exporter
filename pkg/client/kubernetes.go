package client

import (
	"context"

	"golang.org/x/xerrors"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesClient struct {
	Inner kubernetes.Interface
}

func (c *KubernetesClient) Containers() ([]v1.Container, error) {
	// nolint:prealloc
	var containers []v1.Container

	deployments, err := c.Inner.AppsV1().Deployments("").List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		return nil, xerrors.Errorf("could not get deployment: %w", err)
	}
	for _, deployment := range deployments.Items {
		containers = append(containers, deployment.Spec.Template.Spec.Containers...)
	}

	statefulSets, err := c.Inner.AppsV1().StatefulSets("").List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		return nil, xerrors.Errorf("could not get stateful set: %w", err)
	}
	for _, statefulSet := range statefulSets.Items {
		containers = append(containers, statefulSet.Spec.Template.Spec.Containers...)
	}

	daemonSets, err := c.Inner.AppsV1().DaemonSets("").List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		return nil, xerrors.Errorf("could not get daemon set: %w", err)
	}
	for _, daemonSet := range daemonSets.Items {
		containers = append(containers, daemonSet.Spec.Template.Spec.Containers...)
	}

	return containers, nil
}
