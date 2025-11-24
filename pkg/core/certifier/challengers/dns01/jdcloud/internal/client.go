package internal

import (
	"encoding/json"
	"errors"

	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	domainservice "github.com/jdcloud-api/jdcloud-sdk-go/services/domainservice/apis"
)

// This is a partial copy of https://github.com/jdcloud-api/jdcloud-sdk-go/blob/master/services/domainservice/client/DomainserviceClient.go
// to lightweight the vendor packages in the built binary.
type DomainserviceClient struct {
	core.JDCloudClient
}

func NewDomainserviceClient(credential *core.Credential) *DomainserviceClient {
	if credential == nil {
		return nil
	}

	config := core.NewConfig()
	config.SetEndpoint("domainservice.jdcloud-api.com")

	return &DomainserviceClient{
		core.JDCloudClient{
			Credential:  *credential,
			Config:      *config,
			ServiceName: "domainservice",
			Revision:    "2.0.3",
			Logger:      core.NewDummyLogger(),
		},
	}
}

func (c *DomainserviceClient) CreateResourceRecord(request *domainservice.CreateResourceRecordRequest) (*domainservice.CreateResourceRecordResponse, error) {
	if request == nil {
		return nil, errors.New("Request object is nil.")
	}

	resp, err := c.Send(request, c.ServiceName)
	if err != nil {
		return nil, err
	}

	jdResp := &domainservice.CreateResourceRecordResponse{}
	err = json.Unmarshal(resp, jdResp)
	if err != nil {
		c.Logger.Log(core.LogError, "Unmarshal json failed, resp: %s", string(resp))
		return nil, err
	}

	return jdResp, err
}

func (c *DomainserviceClient) DescribeDomains(request *domainservice.DescribeDomainsRequest) (*domainservice.DescribeDomainsResponse, error) {
	if request == nil {
		return nil, errors.New("Request object is nil.")
	}

	resp, err := c.Send(request, c.ServiceName)
	if err != nil {
		return nil, err
	}

	jdResp := &domainservice.DescribeDomainsResponse{}
	err = json.Unmarshal(resp, jdResp)
	if err != nil {
		c.Logger.Log(core.LogError, "Unmarshal json failed, resp: %s", string(resp))
		return nil, err
	}

	return jdResp, err
}

func (c *DomainserviceClient) DeleteResourceRecord(request *domainservice.DeleteResourceRecordRequest) (*domainservice.DeleteResourceRecordResponse, error) {
	if request == nil {
		return nil, errors.New("Request object is nil.")
	}

	resp, err := c.Send(request, c.ServiceName)
	if err != nil {
		return nil, err
	}

	jdResp := &domainservice.DeleteResourceRecordResponse{}
	err = json.Unmarshal(resp, jdResp)
	if err != nil {
		c.Logger.Log(core.LogError, "Unmarshal json failed, resp: %s", string(resp))
		return nil, err
	}

	return jdResp, err
}
