package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rddl-network/shamir-shareholder-service/types"
)

type IShamirShareholderClient interface {
	GetMnemonic(ctx context.Context) (res types.MnemonicBody, err error)
	PostMnemonic(ctx context.Context, mnemonic string) (err error)
}

type ShamirShareholderClient struct {
	baseURL string
	client  *http.Client
}

func NewShamirShareholderClient(baseURL string, client *http.Client) *ShamirShareholderClient {
	if client == nil {
		client = &http.Client{}
	}
	return &ShamirShareholderClient{
		baseURL: baseURL,
		client:  client,
	}
}

func (ssc *ShamirShareholderClient) GetMnemonic(ctx context.Context) (res types.MnemonicBody, err error) {
	err = ssc.doRequest(ctx, http.MethodGet, ssc.baseURL+"/mnemonic", nil, &res)
	return
}

func (ssc *ShamirShareholderClient) PostMnemonic(ctx context.Context, mnemonic string) (err error) {
	requestBody := types.MnemonicBody{
		Mnemonic: mnemonic,
	}
	err = ssc.doRequest(ctx, http.MethodPost, ssc.baseURL+"/mnemonic", requestBody, nil)
	return
}

func (ssc *ShamirShareholderClient) doRequest(ctx context.Context, method, url string, body interface{}, response interface{}) (err error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ssc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &httpError{StatusCode: resp.StatusCode, Msg: strings.Join(resp.Header["Error"], "\n")}
	}

	if response != nil {
		return json.NewDecoder(resp.Body).Decode(response)
	}

	return
}

type httpError struct {
	StatusCode int
	Msg        string
}

func (e *httpError) Error() string {
	return http.StatusText(e.StatusCode) + ": " + e.Msg
}
