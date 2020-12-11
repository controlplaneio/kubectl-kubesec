package cmd

import (
	"github.com/spf13/cobra"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
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
	restConfig, _, err := pluginutils.InitClientAndConfig()
	if err != nil {
		panic(err)
	}
	c := kubernetes.NewForConfigOrDie(restConfig)
	return c
}
