package httpx

import (
	"crypto/tls"
	"time"

	"humpback/common/response"

	reqv3 "github.com/imroc/req/v3"
	"golang.org/x/exp/slog"
)

type HttpXClient interface {
	Get(url string, query map[string]string, header map[string]string, data any, token string) error
	Put(url string, query map[string]string, header map[string]string, body any, data any, token string) error
	Post(url string, query map[string]string, header map[string]string, body any, data any, token string) error
	Delete(url string, query map[string]string, header map[string]string, data any, token string) error
}

type httpxClient struct {
	client *reqv3.Client
}

func NewHttpXClient(tlsConfig *tls.Config) HttpXClient {
	httpC := reqv3.C().SetTimeout(20 * time.Second)
	if tlsConfig != nil {
		httpC.TLSClientConfig = tlsConfig
	} else {
		httpC.TLSClientConfig.InsecureSkipVerify = true
	}
	return &httpxClient{
		client: httpC,
	}
}

func (hx *httpxClient) Get(url string, query map[string]string, header map[string]string, data any, token string) error {
	resp, err := hx.client.R().SetBearerAuthToken(token).SetQueryParams(query).SetHeaders(header).SetSuccessResult(data).Get(url)
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

func (hx *httpxClient) Put(url string, query map[string]string, header map[string]string, body any, data any, token string) error {
	resp, err := hx.client.R().SetBearerAuthToken(token).SetQueryParams(query).SetHeaders(header).SetSuccessResult(data).SetBody(body).Put(url)
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

func (hx *httpxClient) Post(url string, query map[string]string, header map[string]string, body any, data any, token string) error {
	resp, err := hx.client.R().SetBearerAuthToken(token).SetQueryParams(query).SetHeaders(header).SetSuccessResult(data).SetBody(body).Post(url)
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

func (hx *httpxClient) Delete(url string, query map[string]string, header map[string]string, data any, token string) error {
	resp, err := hx.client.R().SetBearerAuthToken(token).SetQueryParams(query).SetHeaders(header).SetSuccessResult(data).Delete(url)
	if err != nil {
		slog.Error("HTTP DELETE request failed", "error", err)
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
