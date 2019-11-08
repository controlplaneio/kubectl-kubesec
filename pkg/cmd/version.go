package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/controlplaneio/kubectl-kubesec/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: "Prints kubesec version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(version.VERSION)
		return nil
	},
}
