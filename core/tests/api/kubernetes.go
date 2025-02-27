package api

import (
	"bytes"
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	//
	// Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// GetDefaultKubernetesClient returns a kubernetes client
func GetDefaultKubernetesClient() *kubernetes.Clientset {
	kubeconfig := GetDefaultKubeConfigPath()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return client
}

func GetDefaultKubeConfigPath() string {
	// get the value of the environment variable KUBE_CONFIG_FILE
	kubeconfig := os.Getenv("KUBE_CONFIG_FILE")
	if kubeconfig == "" {
		// if KUBECONFIG is not set, use the default location
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	return kubeconfig
}

func GetDefaultDynamicClient() dynamic.Interface {
	kubeconfig := GetDefaultKubeConfigPath()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the client
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return client
}

func GetNodeAddress() string {
	// get the address from env variable first: NODE_ADDRESS
	nodeAddress := os.Getenv("NODE_ADDRESS")
	if nodeAddress != "" {
		return nodeAddress
	}

	client := GetDefaultKubernetesClient()
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return nodes.Items[0].Status.Addresses[0].Address
}

func CreateNamespace(client *kubernetes.Clientset, name string) *v1.Namespace {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"laf.dev/testing": "testing",
			},
		},
	}

	result, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	return result
}

func GetObject(namespace string, name string, gvr schema.GroupVersionResource, out interface{}) error {
	client := GetDefaultDynamicClient()
	obj, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, out)
	if err != nil {
		return err
	}
	return nil
}

func KubeUpdateStatus(name, namespace, statusYaml string, gvr schema.GroupVersionResource) error {
	client := GetDefaultDynamicClient()
	user, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	status := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(statusYaml), &status)
	if err != nil {
		return err
	}
	user.Object["status"] = status
	_, err = client.Resource(gvr).Namespace(namespace).UpdateStatus(context.TODO(), user, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func Exec(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return stderr.String(), err
	}
	return stdout.String(), nil
}

func KubeApplyFromTemplate(yaml string, params map[string]string) (string, error) {
	out := os.Expand(yaml, func(k string) string { return params[k] })
	return KubeApply(out)
}

func MustKubeApplyFromTemplate(yaml string, params map[string]string) {
	_, err := KubeApplyFromTemplate(yaml, params)
	if err != nil {
		panic(err)
	}
}

func KubeDeleteFromTemplate(yaml string, params map[string]string) (string, error) {
	out := os.Expand(yaml, func(k string) string { return params[k] })
	return KubeDelete(out)
}

func MustKubeDeleteFromTemplate(yaml string, params map[string]string) {
	_, err := KubeDeleteFromTemplate(yaml, params)
	if err != nil {
		panic(err)
	}
}

func KubeApply(yaml string) (string, error) {
	cmd := `kubectl apply -f - <<EOF` + yaml + `EOF`
	out, err := Exec(cmd)
	if err != nil {
		return out, err
	}

	return out, nil
}

func MustKubeApply(yaml string) {
	_, err := KubeApply(yaml)
	if err != nil {
		panic(err)
	}
}

func KubeDelete(yaml string) (string, error) {
	cmd := `kubectl delete -f - <<EOF` + yaml + `EOF`
	out, err := Exec(cmd)
	if err != nil {
		return out, err
	}

	return out, nil
}

func MustKubeDelete(yaml string) {
	_, err := KubeDelete(yaml)
	if err != nil {
		panic(err)
	}
}

func KubeWaitForDeleted(namespace string, groupResourceName string, timeout string) (string, error) {
	cmd := fmt.Sprintf(`kubectl wait --for=delete --timeout=%s %s -n %s`, timeout, groupResourceName, namespace)
	out, err := Exec(cmd)
	if err != nil {
		return out, err
	}

	return out, nil
}

func MustKubeWaitForDeleted(namespace string, groupResourceName string, timeout string) {
	_, err := KubeWaitForDeleted(namespace, groupResourceName, timeout)
	if err != nil {
		panic(err.Error())
	}
}

func KubeWaitForCondition(namespace string, groupResourceName string, condition string, timeout string) (string, error) {
	cmd := fmt.Sprintf(`kubectl wait --for=condition=%s --timeout=%s %s -n %s`, condition, timeout, groupResourceName, namespace)
	out, err := Exec(cmd)
	if err != nil {
		return out, err
	}

	return out, nil
}

func MustKubeWaitForCondition(namespace string, groupResourceName string, condition string, timeout string) {
	_, err := KubeWaitForCondition(namespace, groupResourceName, condition, timeout)
	if err != nil {
		panic(err.Error())
	}
}

func KubeWaitForReady(namespace string, groupResourceName string, timeout string) (string, error) {
	out, err := KubeWaitForCondition(namespace, groupResourceName, "ready", timeout)
	if err != nil {
		return out, err
	}

	return out, nil
}

func MustKubeWaitForReady(namespace string, groupResourceName string, timeout string) {
	_, err := KubeWaitForReady(namespace, groupResourceName, timeout)
	if err != nil {
		panic(err.Error())
	}
}
