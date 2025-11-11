package internal

import (
	"context"
	"errors"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcteo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

// This is a partial copy of https://github.com/TencentCloud/tencentcloud-sdk-go/blob/master/tencentcloud/teo/v20220901/client.go
// to lightweight the vendor packages in the built binary.
type TeoClient struct {
	common.Client
}

func NewTeoClient(credential common.CredentialIface, region string, clientProfile *profile.ClientProfile) (client *TeoClient, err error) {
	client = &TeoClient{}
	client.Init(region).
		WithCredential(credential).
		WithProfile(clientProfile)
	return
}

func (c *TeoClient) CreateDnsRecord(request *tcteo.CreateDnsRecordRequest) (response *tcteo.CreateDnsRecordResponse, err error) {
	return c.CreateDnsRecordWithContext(context.Background(), request)
}

func (c *TeoClient) CreateDnsRecordWithContext(ctx context.Context, request *tcteo.CreateDnsRecordRequest) (response *tcteo.CreateDnsRecordResponse, err error) {
	if request == nil {
		request = tcteo.NewCreateDnsRecordRequest()
	}
	c.InitBaseRequest(&request.BaseRequest, "teo", tcteo.APIVersion, "CreateDnsRecord")

	if c.GetCredential() == nil {
		return nil, errors.New("CreateDnsRecord require credential")
	}

	request.SetContext(ctx)
	response = tcteo.NewCreateDnsRecordResponse()
	err = c.Send(request, response)
	return
}

func (c *TeoClient) DeleteDnsRecords(request *tcteo.DeleteDnsRecordsRequest) (response *tcteo.DeleteDnsRecordsResponse, err error) {
	return c.DeleteDnsRecordsWithContext(context.Background(), request)
}

func (c *TeoClient) DeleteDnsRecordsWithContext(ctx context.Context, request *tcteo.DeleteDnsRecordsRequest) (response *tcteo.DeleteDnsRecordsResponse, err error) {
	if request == nil {
		request = tcteo.NewDeleteDnsRecordsRequest()
	}
	c.InitBaseRequest(&request.BaseRequest, "teo", tcteo.APIVersion, "DeleteDnsRecords")

	if c.GetCredential() == nil {
		return nil, errors.New("DeleteDnsRecords require credential")
	}

	request.SetContext(ctx)
	response = tcteo.NewDeleteDnsRecordsResponse()
	err = c.Send(request, response)
	return
}

func (c *TeoClient) DescribeDnsRecords(request *tcteo.DescribeDnsRecordsRequest) (response *tcteo.DescribeDnsRecordsResponse, err error) {
	return c.DescribeDnsRecordsWithContext(context.Background(), request)
}

func (c *TeoClient) DescribeDnsRecordsWithContext(ctx context.Context, request *tcteo.DescribeDnsRecordsRequest) (response *tcteo.DescribeDnsRecordsResponse, err error) {
	if request == nil {
		request = tcteo.NewDescribeDnsRecordsRequest()
	}
	c.InitBaseRequest(&request.BaseRequest, "teo", tcteo.APIVersion, "DescribeDnsRecords")

	if c.GetCredential() == nil {
		return nil, errors.New("DescribeDnsRecords require credential")
	}

	request.SetContext(ctx)
	response = tcteo.NewDescribeDnsRecordsResponse()
	err = c.Send(request, response)
	return
}

func (c *TeoClient) ModifyDnsRecordsStatus(request *tcteo.ModifyDnsRecordsStatusRequest) (response *tcteo.ModifyDnsRecordsStatusResponse, err error) {
	return c.ModifyDnsRecordsStatusWithContext(context.Background(), request)
}

func (c *TeoClient) ModifyDnsRecordsStatusWithContext(ctx context.Context, request *tcteo.ModifyDnsRecordsStatusRequest) (response *tcteo.ModifyDnsRecordsStatusResponse, err error) {
	if request == nil {
		request = tcteo.NewModifyDnsRecordsStatusRequest()
	}
	c.InitBaseRequest(&request.BaseRequest, "teo", tcteo.APIVersion, "ModifyDnsRecordsStatus")

	if c.GetCredential() == nil {
		return nil, errors.New("ModifyDnsRecordsStatus require credential")
	}

	request.SetContext(ctx)
	response = tcteo.NewModifyDnsRecordsStatusResponse()
	err = c.Send(request, response)
	return
}
