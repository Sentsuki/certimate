package internal

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	alislb "github.com/alibabacloud-go/slb-20140515/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// This is a partial copy of https://github.com/alibabacloud-go/slb-20140515/blob/master/client/client.go
// to lightweight the vendor packages in the built binary.
type SlbClient struct {
	openapi.Client
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

func (client *SlbClient) DescribeDomainExtensionsWithOptions(request *alislb.DescribeDomainExtensionsRequest, runtime *util.RuntimeOptions) (_result *alislb.DescribeDomainExtensionsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}

	if !tea.BoolValue(util.IsUnset(request.DomainExtensionId)) {
		query["DomainExtensionId"] = request.DomainExtensionId
	}

	if !tea.BoolValue(util.IsUnset(request.ListenerPort)) {
		query["ListenerPort"] = request.ListenerPort
	}

	if !tea.BoolValue(util.IsUnset(request.LoadBalancerId)) {
		query["LoadBalancerId"] = request.LoadBalancerId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDomainExtensions"),
		Version:     tea.String("2014-05-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &alislb.DescribeDomainExtensionsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *SlbClient) DescribeDomainExtensions(request *alislb.DescribeDomainExtensionsRequest) (_result *alislb.DescribeDomainExtensionsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &alislb.DescribeDomainExtensionsResponse{}
	_body, _err := client.DescribeDomainExtensionsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *SlbClient) DescribeLoadBalancerListenersWithOptions(request *alislb.DescribeLoadBalancerListenersRequest, runtime *util.RuntimeOptions) (_result *alislb.DescribeLoadBalancerListenersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}

	if !tea.BoolValue(util.IsUnset(request.Description)) {
		query["Description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.ListenerPort)) {
		query["ListenerPort"] = request.ListenerPort
	}

	if !tea.BoolValue(util.IsUnset(request.ListenerProtocol)) {
		query["ListenerProtocol"] = request.ListenerProtocol
	}

	if !tea.BoolValue(util.IsUnset(request.LoadBalancerId)) {
		query["LoadBalancerId"] = request.LoadBalancerId
	}

	if !tea.BoolValue(util.IsUnset(request.MaxResults)) {
		query["MaxResults"] = request.MaxResults
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["NextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Tag)) {
		query["Tag"] = request.Tag
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeLoadBalancerListeners"),
		Version:     tea.String("2014-05-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &alislb.DescribeLoadBalancerListenersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *SlbClient) DescribeLoadBalancerListeners(request *alislb.DescribeLoadBalancerListenersRequest) (_result *alislb.DescribeLoadBalancerListenersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &alislb.DescribeLoadBalancerListenersResponse{}
	_body, _err := client.DescribeLoadBalancerListenersWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *SlbClient) DescribeLoadBalancerAttributeWithOptions(request *alislb.DescribeLoadBalancerAttributeRequest, runtime *util.RuntimeOptions) (_result *alislb.DescribeLoadBalancerAttributeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}

	if !tea.BoolValue(util.IsUnset(request.LoadBalancerId)) {
		query["LoadBalancerId"] = request.LoadBalancerId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeLoadBalancerAttribute"),
		Version:     tea.String("2014-05-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &alislb.DescribeLoadBalancerAttributeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *SlbClient) DescribeLoadBalancerAttribute(request *alislb.DescribeLoadBalancerAttributeRequest) (_result *alislb.DescribeLoadBalancerAttributeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &alislb.DescribeLoadBalancerAttributeResponse{}
	_body, _err := client.DescribeLoadBalancerAttributeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *SlbClient) DescribeLoadBalancerHTTPSListenerAttributeWithOptions(request *alislb.DescribeLoadBalancerHTTPSListenerAttributeRequest, runtime *util.RuntimeOptions) (_result *alislb.DescribeLoadBalancerHTTPSListenerAttributeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}

	if !tea.BoolValue(util.IsUnset(request.ListenerPort)) {
		query["ListenerPort"] = request.ListenerPort
	}

	if !tea.BoolValue(util.IsUnset(request.LoadBalancerId)) {
		query["LoadBalancerId"] = request.LoadBalancerId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeLoadBalancerHTTPSListenerAttribute"),
		Version:     tea.String("2014-05-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &alislb.DescribeLoadBalancerHTTPSListenerAttributeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *SlbClient) DescribeLoadBalancerHTTPSListenerAttribute(request *alislb.DescribeLoadBalancerHTTPSListenerAttributeRequest) (_result *alislb.DescribeLoadBalancerHTTPSListenerAttributeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &alislb.DescribeLoadBalancerHTTPSListenerAttributeResponse{}
	_body, _err := client.DescribeLoadBalancerHTTPSListenerAttributeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *SlbClient) SetDomainExtensionAttributeWithOptions(request *alislb.SetDomainExtensionAttributeRequest, runtime *util.RuntimeOptions) (_result *alislb.SetDomainExtensionAttributeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}

	if !tea.BoolValue(util.IsUnset(request.DomainExtensionId)) {
		query["DomainExtensionId"] = request.DomainExtensionId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
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

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("SetDomainExtensionAttribute"),
		Version:     tea.String("2014-05-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &alislb.SetDomainExtensionAttributeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *SlbClient) SetDomainExtensionAttribute(request *alislb.SetDomainExtensionAttributeRequest) (_result *alislb.SetDomainExtensionAttributeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &alislb.SetDomainExtensionAttributeResponse{}
	_body, _err := client.SetDomainExtensionAttributeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *SlbClient) SetLoadBalancerHTTPSListenerAttributeWithOptions(request *alislb.SetLoadBalancerHTTPSListenerAttributeRequest, runtime *util.RuntimeOptions) (_result *alislb.SetLoadBalancerHTTPSListenerAttributeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}

	if !tea.BoolValue(util.IsUnset(request.AclId)) {
		query["AclId"] = request.AclId
	}

	if !tea.BoolValue(util.IsUnset(request.AclStatus)) {
		query["AclStatus"] = request.AclStatus
	}

	if !tea.BoolValue(util.IsUnset(request.AclType)) {
		query["AclType"] = request.AclType
	}

	if !tea.BoolValue(util.IsUnset(request.Bandwidth)) {
		query["Bandwidth"] = request.Bandwidth
	}

	if !tea.BoolValue(util.IsUnset(request.CACertificateId)) {
		query["CACertificateId"] = request.CACertificateId
	}

	if !tea.BoolValue(util.IsUnset(request.Cookie)) {
		query["Cookie"] = request.Cookie
	}

	if !tea.BoolValue(util.IsUnset(request.CookieTimeout)) {
		query["CookieTimeout"] = request.CookieTimeout
	}

	if !tea.BoolValue(util.IsUnset(request.Description)) {
		query["Description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.EnableHttp2)) {
		query["EnableHttp2"] = request.EnableHttp2
	}

	if !tea.BoolValue(util.IsUnset(request.Gzip)) {
		query["Gzip"] = request.Gzip
	}

	if !tea.BoolValue(util.IsUnset(request.HealthCheck)) {
		query["HealthCheck"] = request.HealthCheck
	}

	if !tea.BoolValue(util.IsUnset(request.HealthCheckConnectPort)) {
		query["HealthCheckConnectPort"] = request.HealthCheckConnectPort
	}

	if !tea.BoolValue(util.IsUnset(request.HealthCheckDomain)) {
		query["HealthCheckDomain"] = request.HealthCheckDomain
	}

	if !tea.BoolValue(util.IsUnset(request.HealthCheckHttpCode)) {
		query["HealthCheckHttpCode"] = request.HealthCheckHttpCode
	}

	if !tea.BoolValue(util.IsUnset(request.HealthCheckInterval)) {
		query["HealthCheckInterval"] = request.HealthCheckInterval
	}

	if !tea.BoolValue(util.IsUnset(request.HealthCheckMethod)) {
		query["HealthCheckMethod"] = request.HealthCheckMethod
	}

	if !tea.BoolValue(util.IsUnset(request.HealthCheckTimeout)) {
		query["HealthCheckTimeout"] = request.HealthCheckTimeout
	}

	if !tea.BoolValue(util.IsUnset(request.HealthCheckURI)) {
		query["HealthCheckURI"] = request.HealthCheckURI
	}

	if !tea.BoolValue(util.IsUnset(request.HealthyThreshold)) {
		query["HealthyThreshold"] = request.HealthyThreshold
	}

	if !tea.BoolValue(util.IsUnset(request.IdleTimeout)) {
		query["IdleTimeout"] = request.IdleTimeout
	}

	if !tea.BoolValue(util.IsUnset(request.ListenerPort)) {
		query["ListenerPort"] = request.ListenerPort
	}

	if !tea.BoolValue(util.IsUnset(request.LoadBalancerId)) {
		query["LoadBalancerId"] = request.LoadBalancerId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.RequestTimeout)) {
		query["RequestTimeout"] = request.RequestTimeout
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Scheduler)) {
		query["Scheduler"] = request.Scheduler
	}

	if !tea.BoolValue(util.IsUnset(request.ServerCertificateId)) {
		query["ServerCertificateId"] = request.ServerCertificateId
	}

	if !tea.BoolValue(util.IsUnset(request.StickySession)) {
		query["StickySession"] = request.StickySession
	}

	if !tea.BoolValue(util.IsUnset(request.StickySessionType)) {
		query["StickySessionType"] = request.StickySessionType
	}

	if !tea.BoolValue(util.IsUnset(request.TLSCipherPolicy)) {
		query["TLSCipherPolicy"] = request.TLSCipherPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.UnhealthyThreshold)) {
		query["UnhealthyThreshold"] = request.UnhealthyThreshold
	}

	if !tea.BoolValue(util.IsUnset(request.VServerGroup)) {
		query["VServerGroup"] = request.VServerGroup
	}

	if !tea.BoolValue(util.IsUnset(request.VServerGroupId)) {
		query["VServerGroupId"] = request.VServerGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.XForwardedFor)) {
		query["XForwardedFor"] = request.XForwardedFor
	}

	if !tea.BoolValue(util.IsUnset(request.XForwardedFor_ClientSrcPort)) {
		query["XForwardedFor_ClientSrcPort"] = request.XForwardedFor_ClientSrcPort
	}

	if !tea.BoolValue(util.IsUnset(request.XForwardedFor_SLBID)) {
		query["XForwardedFor_SLBID"] = request.XForwardedFor_SLBID
	}

	if !tea.BoolValue(util.IsUnset(request.XForwardedFor_SLBIP)) {
		query["XForwardedFor_SLBIP"] = request.XForwardedFor_SLBIP
	}

	if !tea.BoolValue(util.IsUnset(request.XForwardedFor_SLBPORT)) {
		query["XForwardedFor_SLBPORT"] = request.XForwardedFor_SLBPORT
	}

	if !tea.BoolValue(util.IsUnset(request.XForwardedFor_proto)) {
		query["XForwardedFor_proto"] = request.XForwardedFor_proto
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("SetLoadBalancerHTTPSListenerAttribute"),
		Version:     tea.String("2014-05-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &alislb.SetLoadBalancerHTTPSListenerAttributeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *SlbClient) SetLoadBalancerHTTPSListenerAttribute(request *alislb.SetLoadBalancerHTTPSListenerAttributeRequest) (_result *alislb.SetLoadBalancerHTTPSListenerAttributeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &alislb.SetLoadBalancerHTTPSListenerAttributeResponse{}
	_body, _err := client.SetLoadBalancerHTTPSListenerAttributeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}
