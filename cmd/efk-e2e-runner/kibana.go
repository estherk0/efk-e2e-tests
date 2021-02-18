package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type createRequestBody struct {
	Attributes createAttributes `json:"attributes"`
}

type createAttributes struct {
	Title string `json:"title"`
}

// RunKibanaE2ETest is a func to run e2e tests for Kibana using API.
func RunKibanaE2ETest(kibanaURL string) error {
	err := createKibanaIndexPattern(kibanaURL, "test")
	if err != nil {
		return err
	}

	err = deleteKibanaIndexPattern(kibanaURL, "test")
	if err != nil {
		return err
	}
	return nil
}

func createKibanaIndexPattern(kibanaURL string, title string) error {
	log.Println("[INFO] Kibana URL: " + kibanaURL)
	log.Println("[INFO] Trying to create Kibana Index Pattern \"" + title + "\"...")
	if strings.Contains(kibanaURL, "https") {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	reqBody := &createRequestBody{Attributes: createAttributes{Title: title}}
	reqBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", kibanaURL+"/api/saved_objects/index-pattern/"+title, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("kbn-xsrf", "true")
	req.SetBasicAuth(userid, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	json.Unmarshal(respBody, &res)

	if err != nil {
		return err
	}
	return nil
}

func deleteKibanaIndexPattern(kibanaURL string, title string) error {
	log.Println("[INFO] Trying to delete Kibana Index Pattern \"" + title + "\"...")
	req, err := http.NewRequest("DELETE", kibanaURL+"/api/saved_objects/index-pattern/"+title, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("kbn-xsrf", "true")
	req.SetBasicAuth(userid, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
