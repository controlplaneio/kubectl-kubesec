package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var rootCmd = &cobra.Command{
	Use:   "scan",
	Short: "kubesec.io kubectl plugin",
	Long: `
kubesec.io command line utilities`,
	SilenceUsage: true,
}

var (
	namespace, kubeconfig string
	serializer            *kjson.Serializer
	scanTimeOut           int
	scanURL               string
	kubeContext           string // kubernetes context to use that is present in the kubeconfig
)

func Execute() {
	// global flags
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	rootCmd.PersistentFlags().IntVarP(&scanTimeOut, "timeout", "t", 120, "Scan timeout in seconds")
	rootCmd.PersistentFlags().StringVarP(&scanURL, "url", "u", "https://v2.kubesec.io", "URL to send the request for scanning")
	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig, overrides KUBECONFIG environment variable")
	rootCmd.PersistentFlags().StringVarP(&kubeContext, "context", "c", "", "kubernetes context to use in kubeconfig")

	serializer = kjson.NewSerializerWithOptions(kjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, kjson.SerializerOptions{
		Yaml: true,
	})

	// commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(podCmd)
	rootCmd.AddCommand(deploymentCmd)
	rootCmd.AddCommand(daemonsetCmd)
	rootCmd.AddCommand(statefulsetCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// newKubeClient creates a new Kubernetes client
func newKubeClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := newClientConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to create client config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, fmt.Errorf("Unable to create a new Clienset: %w", err)
	}

	return clientset, nil
}

// newClientConfig creates a new rest client config by fetching
// kubeconfig in the following order:
// 1) CLI flag, 2) KUBECONFIG env var, 3) $HOME/.kube/config
func newClientConfig(kubeconfig string) (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// if empty, it will try to load it from other places
	loadingRules.ExplicitPath = kubeconfig

	configOverrides := &clientcmd.ConfigOverrides{}

	// if kubeContext is not empty, it will override the current context
	if kubeContext != "" {
		configOverrides.CurrentContext = kubeContext
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Unable to create client config: %w", err)
	}

	return config, nil
}
