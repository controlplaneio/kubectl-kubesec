package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	_ "github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/kubectl/pkg/pluginutils"
)

func init() {
	flag.CommandLine.Set("logtostderr", "true")
	flag.CommandLine.Set("v", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_V"))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: kubectl plugin scan DEPLOYMENT_NAME")
		os.Exit(1)
	}

	deploymentName := os.Args[1]
	getDeployment(deploymentName)
}

func loadConfig() (*kubernetes.Clientset, string) {
	restConfig, kubeConfig, err := pluginutils.InitClientAndConfig()
	if err != nil {
		panic(err)
	}
	c := kubernetes.NewForConfigOrDie(restConfig)
	ns, _, _ := kubeConfig.Namespace()
	return c, ns
}

func getDeployment(deploymentName string) {
	client, ns := loadConfig()
	dep, err := client.AppsV1beta2().Deployments(ns).Get(deploymentName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	dep.TypeMeta = metav1.TypeMeta{
		Kind:       "Deployment",
		APIVersion: "apps/v1beta2",
	}

	e := kjson.NewYAMLSerializer(kjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = e.Encode(dep, writer)
	if err != nil {
		panic(err)
	}
	writer.Flush()
	//fmt.Println(b.String())

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", "deployment.yaml")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(fileWriter, &b)
	if err != nil {
		panic(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post("https://kubesec.io/", contentType, bodyBuf)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if len(body) < 1 {
		panic("Unknown result")
	}

	var result KubesecResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	result.print(deploymentName)
}

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(result))
	return nil
}

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

func (r KubesecResult) print(resource string) {
	fmt.Println(fmt.Sprintf("%v kubesec.io score %v", resource, r.Score))
	fmt.Println("-----------------")
	if len(r.Scoring.Critical) > 0 {
		fmt.Println("Critical")
		for i, el := range r.Scoring.Critical {
			fmt.Println(fmt.Sprintf("%v. %v", i+1, el.Selector))
			fmt.Println(el.Reason)

		}
		fmt.Println("-----------------")
	}
	if len(r.Scoring.Advise) > 0 {
		fmt.Println("Advise")
		for i, el := range r.Scoring.Advise {
			fmt.Println(fmt.Sprintf("%v. %v", i+1, el.Selector))
			fmt.Println(el.Reason)
		}
		fmt.Println("-----------------")
	}
}
