package kubesec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
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
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", "object.yaml")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fileWriter, &def)
	if err != nil {
		return nil, err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post("https://kubesec.io/", contentType, bodyBuf)
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

	var result KubesecResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
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
	io.WriteString(w, fmt.Sprintf("kubesec.io score: %v", r.Score))
	io.WriteString(w, "-----------------")
	if len(r.Scoring.Critical) > 0 {
		io.WriteString(w, "Critical")
		for i, el := range r.Scoring.Critical {
			io.WriteString(w, fmt.Sprintf("%v. %v", i+1, el.Selector))
			if len(el.Reason) > 0 {
				io.WriteString(w, el.Reason)
			}

		}
		io.WriteString(w, "-----------------")
	}
	if len(r.Scoring.Advise) > 0 {
		io.WriteString(w, "Advise")
		for i, el := range r.Scoring.Advise {
			io.WriteString(w, fmt.Sprintf("%v. %v", i+1, el.Selector))
			if len(el.Reason) > 0 {
				io.WriteString(w, el.Reason)
			}
		}
	}
}
