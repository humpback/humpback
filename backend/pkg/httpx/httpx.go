package httpx

import (
	"time"

	reqv3 "github.com/imroc/req/v3"
	"humpback/common/response"
)

type HttpXClient interface {
	Get(url string, query map[string]string, header map[string]string, data any) error
	Put(url string, query map[string]string, header map[string]string, body any, data any) error
	Post(url string, query map[string]string, header map[string]string, body any, data any) error
	Delete(url string, query map[string]string, header map[string]string, data any) error
}

type httpxClient struct {
	client *reqv3.Client
}

func NewHttpXClient() HttpXClient {
	httpC := reqv3.C().SetTimeout(20 * time.Second)
	httpC.TLSClientConfig.InsecureSkipVerify = true
	return &httpxClient{
		client: httpC,
	}
}

func (hx *httpxClient) Get(url string, query map[string]string, header map[string]string, data any) error {
	resp, err := hx.client.R().SetQueryParams(query).SetHeaders(header).SetSuccessResult(data).Get(url)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if resp.IsSuccessState() {
		return nil
	}
	respBody, err := resp.ToBytes()
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if resp.GetStatusCode() < 500 {
		errMsg := &response.ErrInfo{}
		if err = resp.UnmarshalJson(errMsg); err != nil {
			return response.NewRespServerErr(err.Error(), string(respBody))
		}
		errMsg.CopyToBizData()
		return errMsg
	}
	return response.NewRespServerErr(string(respBody))

}

func (hx *httpxClient) Put(url string, query map[string]string, header map[string]string, body any, data any) error {
	resp, err := hx.client.R().SetQueryParams(query).SetHeaders(header).SetSuccessResult(data).SetBody(body).Put(url)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}

	if resp.IsSuccessState() {
		return nil
	}
	respBody, err := resp.ToBytes()
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if resp.GetStatusCode() < 500 {
		errMsg := &response.ErrInfo{}
		if err = resp.UnmarshalJson(errMsg); err != nil {
			return response.NewRespServerErr(err.Error(), string(respBody))
		}
		errMsg.CopyToBizData()
		return errMsg
	}
	return response.NewRespServerErr(string(respBody))
}

func (hx *httpxClient) Post(url string, query map[string]string, header map[string]string, body any, data any) error {
	resp, err := hx.client.R().SetQueryParams(query).SetHeaders(header).SetSuccessResult(data).SetBody(body).Post(url)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if resp.IsSuccessState() {
		return nil
	}
	respBody, err := resp.ToBytes()
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if resp.GetStatusCode() < 500 {
		errMsg := &response.ErrInfo{}
		if err = resp.UnmarshalJson(errMsg); err != nil {
			return response.NewRespServerErr(err.Error(), string(respBody))
		}
		errMsg.CopyToBizData()
		return errMsg
	}
	return response.NewRespServerErr(string(respBody))
}

func (hx *httpxClient) Delete(url string, query map[string]string, header map[string]string, data any) error {
	resp, err := hx.client.R().SetQueryParams(query).SetHeaders(header).SetSuccessResult(data).Delete(url)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if resp.IsSuccessState() || resp.GetStatusCode() == 404 {
		return nil
	}
	respBody, err := resp.ToBytes()
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if resp.GetStatusCode() < 500 {
		errMsg := &response.ErrInfo{}
		if err = resp.UnmarshalJson(errMsg); err != nil {
			return response.NewRespServerErr(err.Error(), string(respBody))
		}
		errMsg.CopyToBizData()
		return errMsg
	}
	return response.NewRespServerErr(string(respBody))
}
