package rackspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type BlockDeviceMapping struct {
	BootIndex           int    `json:"boot_index,omitempty"`
	Uuid                int    `json:"uuid,omitempty"`
	BootSourceType      string `json:"boot_source_type:omitempty"`
	BootDestinationType string `json:"boot_destination_type:omitempty"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}

type Personality struct {
	Path    string `json:"path,omitempty"`
	Content string `json:"content,omitempty"`
}

type Network struct {
	UUID string `json:"uuid,omitempty"`
	Port string `json:"port,omitempty"`
}

type CreateServerRequest struct {
	Name               string                 `json:"name"`
	ImageRef           string                 `json:"imageRef"`
	BlockDeviceMapping []BlockDeviceMapping   `json:"block_device_mapping,omitempty"`
	FlavorRef          string                 `json:"flavorRef"`
	ConfigDrive        string                 `json:"config_drive,omitempty"`
	KeyName            string                 `json:"key_name,omitempty"`
	OsDcf              string                 `json:"OS-DCF:diskConfig,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
	Personality        []Personality          `json:"personality,omitempty"`
	UserData           string                 `json:"user_data,omitempty"`
	Networks           []Network              `json:"network,omitempty"`
}

type Link struct {
	Href string `json:"href,omitempty"`
	Rel  string `json:"rel,omitempty"`
}

type CreateServerResponse struct {
	Id           string `json:"id"`
	Links        []Link `json:"link"`
	AdminPass    string `json:"adminPass"`
	OSDiskConfig string `json:"OS- DCF:diskConfig"`
}

func (c RackspaceClient) CreateServer(createServerRequest CreateServerRequest) (CreateServerResponse, error) {
	url := fmt.Sprintf("%s/servers", c.Url)
	fullCreateServerRequest := map[string]CreateServerRequest{
		"server": createServerRequest,
	}
	requestBody, _ := json.Marshal(fullCreateServerRequest)
	log.Printf("[DEBUG] create body: %s", string(requestBody))
	body := bytes.NewReader(requestBody)
	request, _ := http.NewRequest("POST", url, body)
	b, err := c.send(request, []int{http.StatusAccepted})
	if err != nil {
		return CreateServerResponse{}, err
	}
	createServerResponse := map[string]CreateServerResponse{
		"server": CreateServerResponse{},
	}
	json.Unmarshal(b, &createServerResponse)
	return createServerResponse["server"], nil
}

type UpdateServerRequest struct {
}

type UpdateServerResponse struct {
}

func (c RackspaceClient) UpdateServer(updateServerRequest UpdateServerRequest) (UpdateServerResponse, error) {
	return UpdateServerResponse{}, nil
}

func (c RackspaceClient) DestroyServer(id string) error {
	url := fmt.Sprintf("%s/servers/%s", c.Url, id)
	request, _ := http.NewRequest("DELETE", url, nil)
	_, err := c.send(request, []int{204})
	return err
}

type Address struct {
	Addr    string `json:"addr"`
	Version int    `json:"version"`
}

type Flavor struct {
	Id    string `json:"id"`
	Links []Link `json:"links"`
}

type Image struct {
	Id    string `json:"id"`
	Links []Link `json:"links"`
}

type ReadServerResponse struct {
	AccessIPv4    string               `json:"accessIPv4"`
	AccessIPv6    string               `json:"accessIPv6"`
	Addresses     map[string][]Address `json:"addresses"`
	Id            string               `json:"id"`
	Created       string               `json:"created"`
	Flavor        Flavor               `json:"flavour"`
	Image         Image                `json:"image"`
	HostId        string               `json:"hostId"`
	Links         []Link               `json:"links"`
	Metadata      map[string]string    `json:"metadata"`
	Name          string               `json:"name"`
	Progress      int                  `json:"progress"`
	Status        string               `json:"status"`
	TenantId      string               `json:"tennant_id"`
	Updated       string               `json:"updated"`
	UserId        string               `json:"user_id"`
	OsDiskConfig  string               `json:"OS-DCF:diskConfig"`
	ImageSchedule string               `json:"RAX- SI:image_schedule"`
	OsExtSTS      string               `json:"OS-EXT-STS"`
	PublicIpZone  string               `json:"RAX- PUBLIC-IP-ZONE- ID:publicIPZoneId"`
}

func (c RackspaceClient) ReadServer(id string) (ReadServerResponse, error) {
	url := fmt.Sprintf("%s/servers/%s", c.Url, id)
	request, _ := http.NewRequest("GET", url, nil)
	b, err := c.send(request, []int{200, 203, 300})
	if err != nil {
		return ReadServerResponse{}, err
	}
	readServerResponseFull := map[string]ReadServerResponse{
		"server": ReadServerResponse{},
	}
	jsonerr := json.Unmarshal(b, &readServerResponseFull)
	if jsonerr != nil {
		return ReadServerResponse{}, jsonerr
	}
	return readServerResponseFull["server"], nil
}
