package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/kubectl/pkg/pluginutils"
)

func init() {
	flag.CommandLine.Set("logtostderr", "true")
	flag.CommandLine.Set("v", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_V"))
}

const (
	usage   = "usage: kubectl plugin scan [pod|deployment|statefulset|daemonset]/name"
	unknown = "unknown type must be pod, deployment, statefulset or daemonset"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	resource := os.Args[1]
	parts := strings.Split(resource, "/")

	if len(parts) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	client, ns := loadConfig()
	serializer := kjson.NewYAMLSerializer(kjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	switch parts[0] {
	case "pod":
		fmt.Println("scanning pod", parts[1])
		pod, err := client.CoreV1().Pods(ns).Get(parts[1], metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pod.TypeMeta = metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		}
		err = serializer.Encode(pod, writer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "deployment":
		fmt.Println("scanning deployment", parts[1])
		dep, err := client.AppsV1beta2().Deployments(ns).Get(parts[1], metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dep.TypeMeta = metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		}
		err = serializer.Encode(dep, writer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "statefulset":
		fmt.Println("scanning statefulset", parts[1])
		ss, err := client.AppsV1beta2().StatefulSets(ns).Get(parts[1], metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ss.TypeMeta = metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		}
		err = serializer.Encode(ss, writer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "daemonset":
		fmt.Println("scanning daemonset", parts[1])
		ds, err := client.AppsV1beta2().DaemonSets(ns).Get(parts[1], metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ds.TypeMeta = metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		}
		err = serializer.Encode(ds, writer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println(parts[0], unknown)
		os.Exit(1)
	}

	writer.Flush()

	result, err := getResult(buffer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result.print(resource)
}

func loadConfig() (*kubernetes.Clientset, string) {
	restConfig, kubeConfig, err := pluginutils.InitClientAndConfig()
	if err != nil {
		panic(err)
	}
	c := kubernetes.NewForConfigOrDie(restConfig)
	ns, _, _ := kubeConfig.Namespace()
	return c, ns
}
