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

var daemonsetCmd = &cobra.Command{
	Use:     `daemonset [name]`,
	Aliases: []string{"ds"},
	Short:   "Scans daemonset object",
	Example: `  daemonset podinfo --namespace=default`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("daemonset name is required")
		}
		name := args[0]

		fmt.Println("scanning daemonset", name, "in namespace", namespace)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		daemonset, err := kubeClient.AppsV1().DaemonSets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		daemonset.TypeMeta = metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		}
		var buf []byte
		buf, err = json.Marshal(daemonset)
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
