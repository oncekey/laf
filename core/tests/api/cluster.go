package api

import (
	"context"
	"github.com/labring/laf/core/pkg/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteSystemNamespace() {
	name := common.GetSystemNamespace()
	DeleteNamespace(name)
}

func EnsureDefaultSystemNamespace() *v1.Namespace {
	name := common.GetSystemNamespace()
	client := GetDefaultKubernetesClient()
	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return CreateNamespace(client, name)
	}

	return ns
}

func EnsureNamespace(name string) *v1.Namespace {
	client := GetDefaultKubernetesClient()
	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return CreateNamespace(client, name)
	}
	return ns
}

func DeleteNamespace(name string) {
	client := GetDefaultKubernetesClient()
	err := client.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}
