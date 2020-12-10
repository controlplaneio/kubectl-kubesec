package cmd

import (
	"fmt"

	"github.com/controlplaneio/kubectl-kubesec/v2/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: "Prints kubesec version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(version.VERSION)
		return nil
	},
}
