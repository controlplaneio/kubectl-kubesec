package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/controlplaneio/kubectl-kubesec/v2/pkg/kubesec"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var podCmd = &cobra.Command{
	Use:     `pod [name]`,
	Aliases: []string{"po"},
	Short:   "Scans pod object",
	Example: `  pod podinfo --namespace=default`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("pod name is required")
		}
		name := args[0]

		fmt.Println("scanning pod", name, "in namespace", namespace)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		pod, err := kubeClient.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pod.TypeMeta = metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		}
		var buf []byte
		buf, err = json.Marshal(pod)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		result, err := kubesec.NewClient().ScanDefinition(buf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		kubesec.DumpReport(result, os.Stdout)

		return nil
	},
}
