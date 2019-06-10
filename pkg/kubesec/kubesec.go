package kubesec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// KubesecClient represent a client for kubesec.io.
type KubesecClient struct {
}

// NewClient returns a new client for kubesec.io.
func NewClient() *KubesecClient {
	return &KubesecClient{}
}

// ScanDefinition scans the provided resource definition.
func (kc *KubesecClient) ScanDefinition(def bytes.Buffer) (*KubesecResult, error) {
	url := os.Getenv("KUBESEC_URL")
	if url == "" {
		url = "https://v2.kubesec.io/scan"
	}

	resp, err := http.Post(url, "application/yaml", bytes.NewBuffer(def.Bytes()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) < 1 {
		return nil, errors.New("failed to scan definition")
	}

	var result []KubesecResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(string(body))
		return nil, fmt.Errorf("json unmarshal error: %s", err.Error())
	}

	return &result[0], nil
}

// KubesecResult represents a result returned by kubesec.io.
type KubesecResult struct {
	Error   string `json:"error"`
	Score   int    `json:"score"`
	Scoring struct {
		Critical []struct {
			Selector string `json:"selector"`
			Reason   string `json:"reason"`
			Weight   int    `json:"weight"`
		} `json:"critical"`
		Advise []struct {
			Selector string `json:"selector"`
			Reason   string `json:"reason"`
			Href     string `json:"href,omitempty"`
		} `json:"advise"`
	} `json:"scoring"`
}

// Dump writes the result in a human-readable format to the specified writer.
func (r *KubesecResult) Dump(w io.Writer) {
	fmt.Fprintf(w, "kubesec.io score: %v\n", r.Score)
	fmt.Fprintln(w, "-----------------")
	if len(r.Scoring.Critical) > 0 {
		fmt.Fprintln(w, "Critical")
		for i, el := range r.Scoring.Critical {
			fmt.Fprintf(w, "%v. %v\n", i+1, el.Selector)
			if len(el.Reason) > 0 {
				fmt.Fprintln(w, el.Reason)
			}

		}
		fmt.Fprintln(w, "-----------------")
	}
	if len(r.Scoring.Advise) > 0 {
		fmt.Fprintln(w, "Advise")
		for i, el := range r.Scoring.Advise {
			fmt.Fprintf(w, "%v. %v\n", i+1, el.Selector)
			if len(el.Reason) > 0 {
				fmt.Fprintln(w, el.Reason)
			}
		}
	}
}
