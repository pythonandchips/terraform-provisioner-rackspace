package rackspace

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateServerRequest(t *testing.T) {
	response := map[string]interface{}{
		"server": map[string]interface{}{
			"id": "123456789",
			"links": []map[string]string{
				{"href": "http://rackspace/0000001",
					"rel": "self"},
				{"href": "http://rackspace/0000002",
					"rel": "bookmark"},
			},
			"adminPass": "password",
		},
	}
	var request *http.Request
	var receivedCreateServerRequest CreateServerRequest
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		body, _ := ioutil.ReadAll(r.Body)
		createServerRequestFull := map[string]CreateServerRequest{
			"server": CreateServerRequest{},
		}
		json.Unmarshal(body, &createServerRequestFull)
		receivedCreateServerRequest = createServerRequestFull["server"]
		b, _ := json.Marshal(response)
		w.WriteHeader(http.StatusAccepted)
		w.Write(b)
	}))
	defer ts.Close()

	c := RackspaceClient{Token: "000111", Url: ts.URL}

	createServerRequest := CreateServerRequest{
		Name:      "foo",
		ImageRef:  "imageRef",
		FlavorRef: "debian",
		Networks: []Network{
			{UUID: "1030202", Port: "8080"},
		},
	}
	createServerResponse, err := c.CreateServer(createServerRequest)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
	if createServerResponse.Id != "123456789" {
		t.Errorf("Expected Id to be %s but was %s", "123456789", createServerResponse.Id)
	}
	if createServerResponse.AdminPass != "password" {
		t.Errorf("Expected admin password to be %s but was %s", "password", createServerResponse.AdminPass)
	}
	if request.Header["X-Auth-Token"][0] != "000111" {
		t.Errorf("Expected request header X-Auth-Token to be %s but was %s", "000111", request.Header["X-Auth-Token"][0])
	}
	if request.Header["Content-Type"][0] != "application/json" {
		t.Errorf("Expected request header Content-Type to be %s but was %s", "000111", request.Header["X-Auth-Token"][0])
	}
	if request.Header["Accept"][0] != "application/json" {
		t.Errorf("Expected request header Accpet to be %s but was %s", "000111", request.Header["X-Auth-Token"][0])
	}
	if receivedCreateServerRequest.Name != "foo" {
		t.Errorf("expected server to recieve name %s but was %s", "foo", receivedCreateServerRequest.Name)
	}
}

func TestReadServer(t *testing.T) {
	response := exampleReadJson
	var request *http.Request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte(response))
	}))
	defer ts.Close()
	c := RackspaceClient{Token: "000111", Url: ts.URL}
	readServerResponse, err := c.ReadServer("9193829")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
	if len(readServerResponse.Addresses["public"]) != 2 {
		t.Errorf("Expected 1 address but got %d", len(readServerResponse.Addresses["internal"]))
	}
	if readServerResponse.Status != "BUILD" {
		t.Errorf("Expected status to be %s but got %s", "BUILD", readServerResponse.Status)
	}
}

var exampleReadJson = `
{
  "server":{
    "status":"BUILD",
    "updated":"2016-05-26T21:07:39Z",
    "hostId":"3e67fc0bf7766cc9806c480d7a892d2318c52527b3bc1acf8d4125fd",
    "addresses":{
      "public":[
        {
          "version":4,
          "addr":"162.13.6.40"
        },
        {
          "version":6,
          "addr":"2a00:1a48:7805:112:be76:4eff:fe08:f513"
        }
      ],
      "private":[
        {
          "version":4,
          "addr":"10.179.2.133"
        }
      ]
    },
    "links":[
      {
        "href":"https://lon.servers.api.rackspacecloud.com/v2/10004527/servers/6081544a-c512-4ac7-bcf5-625362131a64",
        "rel":"self"
      },
      {
        "href":"https://lon.servers.api.rackspacecloud.com/10004527/servers/6081544a-c512-4ac7-bcf5-625362131a64",
        "rel":"bookmark"
      }
    ],
    "key_name":"colin-gemmell-linux",
    "image":{
      "id":"cf16c435-7bed-4dc3-b76e-57b09987866d",
      "links":[
        {
          "href":"https://lon.servers.api.rackspacecloud.com/10004527/images/cf16c435-7bed-4dc3-b76e-57b09987866d",
          "rel":"bookmark"
        }
      ]
    },
    "RAX-PUBLIC-IP-ZONE-ID:publicIPZoneId":"8166bf0918e7e477a90870e5827e3ade84158ecb0726ee102c8094b4",
    "OS-EXT-STS:task_state":"spawning",
    "OS-EXT-STS:vm_state":"building",
    "flavor":{
      "id":"2",
      "links":[
        {
          "href":"https://lon.servers.api.rackspacecloud.com/10004527/flavors/2",
          "rel":"bookmark"
        }
      ]
    },
    "id":"6081544a-c512-4ac7-bcf5-625362131a64",
    "user_id":"10015904",
    "name":"rundeck-01",
    "created":"2016-05-26T21:07:27Z",
    "tenant_id":"10004527",
    "OS-DCF:diskConfig":"MANUAL",
    "accessIPv4":"",
    "accessIPv6":"",
    "progress":70,
    "OS-EXT-STS:power_state":0
  }
}`
