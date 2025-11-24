package internal

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	alislb "github.com/alibabacloud-go/slb-20140515/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/alibabacloud-go/tea/tea"
)

// This is a partial copy of https://github.com/alibabacloud-go/slb-20140515/blob/master/client/client.go
// to lightweight the vendor packages in the built binary.
type SlbClient struct {
	openapi.Client
	DisableSDKError *bool
}

func NewSlbClient(config *openapi.Config) (*SlbClient, error) {
	client := new(SlbClient)
	err := client.Init(config)
	return client, err
}

func (client *SlbClient) Init(config *openapi.Config) (_err error) {
	_err = client.Client.Init(config)
	if _err != nil {
		return _err
	}
	_err = client.CheckConfig(config)
	if _err != nil {
		return _err
	}

	return nil
}

func (client *SlbClient) DescribeServerCertificatesWithOptions(request *alislb.DescribeServerCertificatesRequest, runtime *util.RuntimeOptions) (_result *alislb.DescribeServerCertificatesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}

	if !dara.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		query["ResourceGroupId"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ServerCertificateId)) {
		query["ServerCertificateId"] = request.ServerCertificateId
	}

	if !tea.BoolValue(util.IsUnset(request.Tag)) {
		query["Tag"] = request.Tag
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeServerCertificates"),
		Version:     tea.String("2014-05-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &alislb.DescribeServerCertificatesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *SlbClient) DescribeServerCertificates(request *alislb.DescribeServerCertificatesRequest) (_result *alislb.DescribeServerCertificatesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &alislb.DescribeServerCertificatesResponse{}
	_body, _err := client.DescribeServerCertificatesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *SlbClient) UploadServerCertificateWithOptions(request *alislb.UploadServerCertificateRequest, runtime *util.RuntimeOptions) (_result *alislb.UploadServerCertificateResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}

	if !tea.BoolValue(util.IsUnset(request.AliCloudCertificateId)) {
		query["AliCloudCertificateId"] = request.AliCloudCertificateId
	}

	if !tea.BoolValue(util.IsUnset(request.AliCloudCertificateName)) {
		query["AliCloudCertificateName"] = request.AliCloudCertificateName
	}

	if !tea.BoolValue(util.IsUnset(request.AliCloudCertificateRegionId)) {
		query["AliCloudCertificateRegionId"] = request.AliCloudCertificateRegionId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PrivateKey)) {
		query["PrivateKey"] = request.PrivateKey
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		query["ResourceGroupId"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ServerCertificate)) {
		query["ServerCertificate"] = request.ServerCertificate
	}

	if !tea.BoolValue(util.IsUnset(request.ServerCertificateName)) {
		query["ServerCertificateName"] = request.ServerCertificateName
	}

	if !tea.BoolValue(util.IsUnset(request.Tag)) {
		query["Tag"] = request.Tag
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("UploadServerCertificate"),
		Version:     tea.String("2014-05-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &alislb.UploadServerCertificateResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *SlbClient) UploadServerCertificate(request *alislb.UploadServerCertificateRequest) (_result *alislb.UploadServerCertificateResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &alislb.UploadServerCertificateResponse{}
	_body, _err := client.UploadServerCertificateWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}
