package internal

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	aliesa "github.com/alibabacloud-go/esa-20240910/v2/client"
	"github.com/alibabacloud-go/tea/dara"
)

// This is a partial copy of https://github.com/alibabacloud-go/esa-20240910/blob/master/client/client.go
// to lightweight the vendor packages in the built binary.
type EsaClient struct {
	openapi.Client
	DisableSDKError *bool
}

func NewEsaClient(config *openapiutil.Config) (*EsaClient, error) {
	client := new(EsaClient)
	err := client.Init(config)
	return client, err
}

func (client *EsaClient) Init(config *openapiutil.Config) (_err error) {
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

func (client *EsaClient) CreateRecordWithOptions(tmpReq *aliesa.CreateRecordRequest, runtime *dara.RuntimeOptions) (_result *aliesa.CreateRecordResponse, _err error) {
	_err = tmpReq.Validate()
	if _err != nil {
		return _result, _err
	}

	request := &aliesa.CreateRecordShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !dara.IsNil(tmpReq.AuthConf) {
		request.AuthConfShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.AuthConf, dara.String("AuthConf"), dara.String("json"))
	}

	if !dara.IsNil(tmpReq.Data) {
		request.DataShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.Data, dara.String("Data"), dara.String("json"))
	}

	query := map[string]interface{}{}
	if !dara.IsNil(request.AuthConfShrink) {
		query["AuthConf"] = request.AuthConfShrink
	}

	if !dara.IsNil(request.BizName) {
		query["BizName"] = request.BizName
	}

	if !dara.IsNil(request.Comment) {
		query["Comment"] = request.Comment
	}

	if !dara.IsNil(request.DataShrink) {
		query["Data"] = request.DataShrink
	}

	if !dara.IsNil(request.HostPolicy) {
		query["HostPolicy"] = request.HostPolicy
	}

	if !dara.IsNil(request.Proxied) {
		query["Proxied"] = request.Proxied
	}

	if !dara.IsNil(request.RecordName) {
		query["RecordName"] = request.RecordName
	}

	if !dara.IsNil(request.SiteId) {
		query["SiteId"] = request.SiteId
	}

	if !dara.IsNil(request.SourceType) {
		query["SourceType"] = request.SourceType
	}

	if !dara.IsNil(request.Ttl) {
		query["Ttl"] = request.Ttl
	}

	if !dara.IsNil(request.Type) {
		query["Type"] = request.Type
	}

	req := &openapiutil.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapiutil.Params{
		Action:      dara.String("CreateRecord"),
		Version:     dara.String("2024-09-10"),
		Protocol:    dara.String("HTTPS"),
		Pathname:    dara.String("/"),
		Method:      dara.String("POST"),
		AuthType:    dara.String("AK"),
		Style:       dara.String("RPC"),
		ReqBodyType: dara.String("formData"),
		BodyType:    dara.String("json"),
	}
	_result = &aliesa.CreateRecordResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = dara.Convert(_body, &_result)
	return _result, _err
}

func (client *EsaClient) CreateRecord(request *aliesa.CreateRecordRequest) (_result *aliesa.CreateRecordResponse, _err error) {
	runtime := &dara.RuntimeOptions{}
	_result = &aliesa.CreateRecordResponse{}
	_body, _err := client.CreateRecordWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *EsaClient) DeleteRecordWithOptions(request *aliesa.DeleteRecordRequest, runtime *dara.RuntimeOptions) (_result *aliesa.DeleteRecordResponse, _err error) {
	_err = request.Validate()
	if _err != nil {
		return _result, _err
	}

	query := map[string]interface{}{}
	if !dara.IsNil(request.RecordId) {
		query["RecordId"] = request.RecordId
	}

	if !dara.IsNil(request.SecurityToken) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapiutil.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapiutil.Params{
		Action:      dara.String("DeleteRecord"),
		Version:     dara.String("2024-09-10"),
		Protocol:    dara.String("HTTPS"),
		Pathname:    dara.String("/"),
		Method:      dara.String("POST"),
		AuthType:    dara.String("AK"),
		Style:       dara.String("RPC"),
		ReqBodyType: dara.String("formData"),
		BodyType:    dara.String("json"),
	}
	_result = &aliesa.DeleteRecordResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = dara.Convert(_body, &_result)
	return _result, _err
}

func (client *EsaClient) DeleteRecord(request *aliesa.DeleteRecordRequest) (_result *aliesa.DeleteRecordResponse, _err error) {
	runtime := &dara.RuntimeOptions{}
	_result = &aliesa.DeleteRecordResponse{}
	_body, _err := client.DeleteRecordWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *EsaClient) ListSitesWithOptions(tmpReq *aliesa.ListSitesRequest, runtime *dara.RuntimeOptions) (_result *aliesa.ListSitesResponse, _err error) {
	_err = tmpReq.Validate()
	if _err != nil {
		return _result, _err
	}

	request := &aliesa.ListSitesShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !dara.IsNil(tmpReq.TagFilter) {
		request.TagFilterShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.TagFilter, dara.String("TagFilter"), dara.String("json"))
	}

	query := openapiutil.Query(dara.ToMap(request))
	req := &openapiutil.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapiutil.Params{
		Action:      dara.String("ListSites"),
		Version:     dara.String("2024-09-10"),
		Protocol:    dara.String("HTTPS"),
		Pathname:    dara.String("/"),
		Method:      dara.String("GET"),
		AuthType:    dara.String("AK"),
		Style:       dara.String("RPC"),
		ReqBodyType: dara.String("formData"),
		BodyType:    dara.String("json"),
	}
	_result = &aliesa.ListSitesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = dara.Convert(_body, &_result)
	return _result, _err
}

func (client *EsaClient) ListSites(request *aliesa.ListSitesRequest) (_result *aliesa.ListSitesResponse, _err error) {
	runtime := &dara.RuntimeOptions{}
	_result = &aliesa.ListSitesResponse{}
	_body, _err := client.ListSitesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}
