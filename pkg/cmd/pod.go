package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/controlplaneio/kubectl-kubesec/pkg/kubesec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

var podCmd = &cobra.Command{
	Use:     `pod [name]`,
	Short:   "Scans pod object",
	Example: `  pod podinfo --namespace=default`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("pod name is required")
		}
		name := args[0]

		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)

		fmt.Println("scanning pod", name, "in namespace", namespace)
		pod, err := kubeClient.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
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

		writer.Flush()

		result, err := kubesec.NewClient().ScanDefinition(buffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if result.Error != "" {
			fmt.Println(result.Error)
			os.Exit(1)
		}

		result.Dump(os.Stdout)

		return nil
	},
}
