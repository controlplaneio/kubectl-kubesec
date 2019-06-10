package main

import (
	"flag"
	"fmt"
	"github.com/controlplaneio/kubectl-kubesec/pkg/cmd"
	_ "github.com/golang/glog"
	"os"
	"strings"
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
