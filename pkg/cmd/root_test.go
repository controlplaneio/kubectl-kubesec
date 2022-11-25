package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewClient tests the creation of the Kubeconfig from
// different location with the overrides.
func TestNewClient(t *testing.T) {
	var (
		// The URL, is used to determine if the right kubeconfig is loaded
		kubeconfigTemplate = `apiVersion: v1
clusters:
- cluster:
    server: _URL_
  name: kind-kind
contexts:
- context:
    cluster: kind-kind
    user: kind-kind
  name: kind-kind
current-context: kind-kind
kind: Config
preferences: {}
users:
- name: kind-kind
`

		kubeconfigFromUser       = strings.Replace(kubeconfigTemplate, "_URL_", "kubeconfig-from-user", 1)
		kubeconfigFromEnv        = strings.Replace(kubeconfigTemplate, "_URL_", "kubeconfig-from-env", 1)
		kubeconfigFromDefaultLoc = strings.Replace(kubeconfigTemplate, "_URL_", "kubeconfig-from-default-loc", 1)
	)

	var tests = []struct {
		name              string
		setFromUser       bool
		setFromEnv        bool
		SetFromDefaultLoc bool
		expectedHost      string
	}{
		{
			name:         "Kubeconfig is provided from CLI",
			setFromUser:  true,
			expectedHost: "kubeconfig-from-user",
		},
		{
			name:         "Kubeconfig is provided on with env var KUBECONFIG",
			setFromEnv:   true,
			expectedHost: "kubeconfig-from-env",
		},
		{
			name:              "Kubeconfig is set from the default loc, not provided from CLI nor env",
			SetFromDefaultLoc: true,
			expectedHost:      "kubeconfig-from-default-loc",
		},
		{
			name:         "Kubeconfig is provided from CLI and env",
			setFromUser:  true,
			setFromEnv:   true,
			expectedHost: "kubeconfig-from-user",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			kubeconfigPath := ""
			dir := t.TempDir()

			switch {
			case tt.setFromUser:
				kubeconfigPath = filepath.Join(dir, "config")
				err := os.WriteFile(kubeconfigPath, []byte(kubeconfigFromUser), 0400)
				if err != nil {
					t.Error(err)
				}
			case tt.setFromEnv:
				kubeconfigPath = filepath.Join(dir, "config")
				err := os.WriteFile(kubeconfigPath, []byte(kubeconfigFromEnv), 0400)
				if err != nil {
					t.Error(err)
				}

				t.Setenv("KUBECONFIG", kubeconfigPath)
			case tt.SetFromDefaultLoc:
				home := dir
				t.Setenv("HOME", home)
				err := os.Mkdir(filepath.Join(home, ".kube"), 0700)
				if err != nil {
					t.Error(err)
				}
				kubeconfigPath = filepath.Join(home, ".kube/config")

				err = os.WriteFile(kubeconfigPath, []byte(kubeconfigFromDefaultLoc), 0400)
				if err != nil {
					t.Error(err)
				}
			}

			config, err := newClientConfig(kubeconfigPath)
			if err != nil {
				t.Error(err)
			}

			if config.Host != tt.expectedHost {
				t.Errorf("Expected Kubeconfig host (location), got: %s, want: %s", config.Host, tt.expectedHost)
			}

			client, err := newKubeClient(kubeconfigPath)
			if err != nil {
				t.Error(err)
			}

			if client == nil {
				t.Errorf("client should not be nil")
			}
		})
	}
}
