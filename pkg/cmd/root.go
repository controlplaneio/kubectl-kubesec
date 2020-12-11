package cmd

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

var RootCmd = &cobra.Command{
	Use:   "scan",
	Short: "kubesec.io kubectl plugin",
	Long: `
kubesec.io command line utilities`,
}

var (
	namespace  string
	kubeClient *kubernetes.Clientset
	serializer *kjson.Serializer
)

func init() {
	// global flags
	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")

	// client
	kubeClient = loadConfig()
	serializer = kjson.NewYAMLSerializer(kjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	// commands
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(podCmd)
	RootCmd.AddCommand(deploymentCmd)
	RootCmd.AddCommand(daemonsetCmd)
	RootCmd.AddCommand(statefulsetCmd)
}

func loadConfig() *kubernetes.Clientset {
	home := os.Getenv("HOME")
	if len(home) == 0 {
		panic("no home")
	}
	clientSet, err := clientSet(home + "/.kube/config")
	if err != nil {
		panic(err)
	}
	return clientSet
}

func clientSet(kubeconfig string) (*kubernetes.Clientset, error) {
	kubeconfigBytes, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		return nil, err
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
