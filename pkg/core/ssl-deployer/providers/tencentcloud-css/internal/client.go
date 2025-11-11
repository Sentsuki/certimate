package internal

import (
	"context"
	"errors"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tclive "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
)

// This is a partial copy of https://github.com/TencentCloud/tencentcloud-sdk-go/blob/master/tencentcloud/live/v20180801/client.go
// to lightweight the vendor packages in the built binary.
type LiveClient struct {
	common.Client
}

func NewLiveClient(credential common.CredentialIface, region string, clientProfile *profile.ClientProfile) (client *LiveClient, err error) {
	client = &LiveClient{}
	client.Init(region).
		WithCredential(credential).
		WithProfile(clientProfile)
	return
}

func (c *LiveClient) ModifyLiveDomainCertBindings(request *tclive.ModifyLiveDomainCertBindingsRequest) (response *tclive.ModifyLiveDomainCertBindingsResponse, err error) {
	return c.ModifyLiveDomainCertBindingsWithContext(context.Background(), request)
}

func (c *LiveClient) ModifyLiveDomainCertBindingsWithContext(ctx context.Context, request *tclive.ModifyLiveDomainCertBindingsRequest) (response *tclive.ModifyLiveDomainCertBindingsResponse, err error) {
	if request == nil {
		request = tclive.NewModifyLiveDomainCertBindingsRequest()
	}
	c.InitBaseRequest(&request.BaseRequest, "live", tclive.APIVersion, "ModifyLiveDomainCertBindings")

	if c.GetCredential() == nil {
		return nil, errors.New("ModifyLiveDomainCertBindings require credential")
	}

	request.SetContext(ctx)
	response = tclive.NewModifyLiveDomainCertBindingsResponse()
	err = c.Send(request, response)
	return
}
