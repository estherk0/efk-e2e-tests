package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	userid   = "elastic"
	password = "tacoword"
)

// ElasticsearchIndex is an index struct from Elasticsearch API
type ElasticsearchIndex struct {
	Health string `json:"health"`
	Index  string `json:"index"`
	Status string `json:"status"`
}

// QueryResponse is a response body from GET /index/search API
type QueryResponse struct {
	Result queryHit `json:"hits"`
}

type queryHit struct {
	Hits []elasticsearchDocument `json:"hits"`
}

type elasticsearchDocument struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	ID     string                 `json:"_id"`
	Source map[string]interface{} `json:"_source"`
}

// RunElasticsearchE2ETest is a func to run e2e tests for Elasticsearch using API.
func RunElasticsearchE2ETest(esURL string) error {
	// check the indicies
	indices, err := getElasticsearchIndices(esURL)
	if err != nil {
		return err
	}

	for _, idx := range indices {
		err = queryAllWithIndex(esURL, idx)
		if err != nil {
			return err
		}
	}
	return nil
}

func getElasticsearchIndices(esURL string) ([]string, error) {
	targetIndices := [1]string{"platform"} // TODO
	log.Println("[INFO] Trying to check Elasticsearch Indices...")
	log.Println("[INFO] Elasticsearch URL: " + esURL)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", esURL+"/_cat/indices", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(userid, password)
	query := req.URL.Query()
	query.Add("format", "JSON")
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.Status != "200 OK" {
		return nil, err
	}
	defer resp.Body.Close()

	var resIndices []ElasticsearchIndex
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(respBody, &resIndices)

	var idxNames = []string{}
	for _, substr := range targetIndices {
		name, err := findOneIndex(resIndices, substr)
		idxNames = append(idxNames, name)
		if err != nil {
			return nil, err
		}
	}
	return idxNames, nil
}

func findOneIndex(indices []ElasticsearchIndex, substr string) (string, error) {
	for _, idx := range indices {
		if strings.Contains(idx.Index, substr) {
			return idx.Index, nil
		}
	}
	return "", errors.New("[ERROR] no matching index for " + substr)
}

func queryAllWithIndex(esURL string, indexName string) error {
	log.Println("[INFO] Trying to search hits with Index \"" + indexName + "\"...")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", esURL+"/"+indexName+"/_search", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(userid, password)
	query := req.URL.Query()
	query.Add("size", "1")
	query.Add("q", "*")
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var queryResp QueryResponse
	json.Unmarshal(respBody, &queryResp)
	if len(queryResp.Result.Hits) == 0 {
		return errors.New("[ERROR] No document hits for the index " + indexName)
	}
	log.Println("[INFO] Found docuemnt hits!")
	log.Printf("%s\n", queryResp.Result.Hits[0].Source)

	return nil
}
