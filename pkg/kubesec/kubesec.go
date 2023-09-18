package kubesec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// KubesecClient represent a client for kubesec.io
type KubesecClient struct {
	URL        string // URL to send the request for scanning
	TimeOutSec int    // Scan timeout in seconds
}

// NewClient returns a new client for kubesec.io.
func NewClient(url string, timeOutSec int) *KubesecClient {
	return &KubesecClient{
		URL:        url,
		TimeOutSec: timeOutSec,
	}
}

// ScanDefinition scans the provided resource definition.
func (kc *KubesecClient) ScanDefinition(def bytes.Buffer) (KubeSecResults, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Duration(kc.TimeOutSec)*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, kc.URL, &def)
	if err != nil {
		return nil, err
	}

	contentType := "application/x-www-form-urlencoded"

	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %v response from %v instead of 200 OK", resp.StatusCode, kc.URL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) < 1 {
		return nil, errors.New("failed to scan definition")
	}

	// API version v2 of Kubesec available at https://v2.kubesec.io returns a slice of results
	var results []KubesecResult

	err = json.Unmarshal(body, &results)

	if err != nil {
		return nil, err
	}

	return results, nil
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
		fmt.Fprintf(w, "Advise")
		for i, el := range r.Scoring.Advise {
			fmt.Fprintf(w, "%v. %v\n", i+1, el.Selector)
			if len(el.Reason) > 0 {
				fmt.Fprintln(w, el.Reason)
			}
		}
	}
}

// KubeSecResults - holds a slice of scan results
type KubeSecResults []KubesecResult

// Dump - calls upstream Dump function and returns an error if any scan object has non-empty error field
func (r KubeSecResults) Dump(w io.Writer) error {
	var msg string
	for _, result := range r {
		if result.Error != "" {
			msg = result.Error + "," + msg
		}
		result.Dump(w)
	}
	if msg != "" {
		return errors.New(strings.TrimSpace(msg))
	}
	return nil
}
