package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/controlplaneio/kubectl-kubesec/v2/pkg/cmd"
	_ "github.com/golang/glog"
)

func init() {
	flag.CommandLine.Set("logtostderr", "true")
}

func main() {
	cmd.RootCmd.SetArgs(os.Args[1:])
	if err := cmd.RootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}
