package rackspace

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client interface {
}

type RackspaceClient struct {
	Token string
	Url   string
}

func NewRackspaceClient(token, tenantId, region string) RackspaceClient {
	url := fmt.Sprintf("https://%s.servers.api.rackspacecloud.com/v2/%s", region, tenantId)
	return RackspaceClient{Token: token, Url: url}
}

func (c RackspaceClient) send(request *http.Request, successCodes []int) ([]byte, error) {
	request.Header = map[string][]string{
		"X-Auth-Token": {c.Token},
		"Content-Type": {"application/json"},
		"Accept":       {"application/json"},
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	isSuccess := false
	for _, code := range successCodes {
		if resp.StatusCode == code {
			isSuccess = true
			break
		}
	}
	if !isSuccess {
		return []byte{}, errors.New(fmt.Sprintf("Create has failed: %s", string(body)))
	}
	log.Printf("[DEBUG] server response: %s", string(body))
	return body, err
}
