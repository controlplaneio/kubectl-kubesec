package main

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

type KubesecResult struct {
	Score   int `json:"score"`
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

func (r *KubesecResult) print(resource string) {
	fmt.Println(fmt.Sprintf("%v kubesec.io score %v", resource, r.Score))
	fmt.Println("-----------------")
	if len(r.Scoring.Critical) > 0 {
		fmt.Println("Critical")
		for i, el := range r.Scoring.Critical {
			fmt.Println(fmt.Sprintf("%v. %v", i+1, el.Selector))
			if len(el.Reason) > 0 {
				fmt.Println(el.Reason)
			}

		}
		fmt.Println("-----------------")
	}
	if len(r.Scoring.Advise) > 0 {
		fmt.Println("Advise")
		for i, el := range r.Scoring.Advise {
			fmt.Println(fmt.Sprintf("%v. %v", i+1, el.Selector))
			if len(el.Reason) > 0 {
				fmt.Println(el.Reason)
			}
		}
	}
}

func getResult(definition bytes.Buffer) (*KubesecResult, error) {

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", "object.yaml")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fileWriter, &definition)
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
		return nil, errors.New("unknown result")
	}

	var result KubesecResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
