package xinnet

type apiResponse interface {
	GetCode() string
	GetMessage() string
	GetRequestId() string
}

type apiResponseBase struct {
	Code      *string `json:"code,omitempty"`
	Message   *string `json:"message,omitempty"`
	RequestId *string `json:"requestId,omitempty"`
}

func (r *apiResponseBase) GetCode() string {
	if r.Code == nil {
		return ""
	}

	return *r.Code
}

func (r *apiResponseBase) GetMessage() string {
	if r.Message == nil {
		return ""
	}

	return *r.Message
}

func (r *apiResponseBase) GetRequestId() string {
	if r.RequestId == nil {
		return ""
	}

	return *r.RequestId
}

var _ apiResponse = (*apiResponseBase)(nil)
