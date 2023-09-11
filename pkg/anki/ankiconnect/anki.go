package ankiconnect

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"
)

const apiVersion = 6

type Anki struct {
	client http.Client
	url    string
	apiKey string
}

type Options struct {
	URL    string
	APIKey string
	// Transport is used for http client.
	// Please not that anki-connect can not properly reuse connections!
	Transport http.RoundTripper
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}

func New(o *Options) (*Anki, error) {
	transport := o.Transport
	if transport == nil {
		transport = &http.Transport{
			// we need to disable keepAlive to deal with anki-connect connection reuse behavior
			DisableKeepAlives: true,
			Proxy:             http.ProxyFromEnvironment,
			DialContext: defaultTransportDialContext(&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}),
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}
	return &Anki{
		client: http.Client{
			Transport: transport,
		},
		apiKey: o.APIKey,
		url:    o.URL,
	}, nil
}

func (a *Anki) request(ctx context.Context, action string, params any, result any) error {
	requestBody := fullRequest{
		Action:  action,
		Version: apiVersion,
		Params:  params,
		Key:     a.apiKey,
	}
	body, err := json.Marshal(&requestBody)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", a.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application/json; charset=utf-8")

	response, err := a.client.Do(request)
	if err != nil {
		return &ConnectionError{Err: err}
	}
	defer response.Body.Close()

	// I saw people do something status < OK || status >= BadRequest.
	// But Go client automatically tries to manage redirects and other
	// reponses from anki connect don't seem to be very useful, at least
	// author doesn't mention any at all.
	if response.StatusCode != http.StatusOK {
		var errResp errResponse
		// we will try to get error (despite anki-connect seems to return OK with error)
		decodeErr := json.NewDecoder(response.Body).Decode(&errResp)
		if decodeErr != nil {
			return newUnableDecodedError(response.StatusCode, decodeErr)
		}
		if errResp.Error == "" {
			return newUnexpectedStatusError(response.StatusCode)
		}
		return newServerError(errResp.Error)
	}
	fullResp := fullResponse{
		Result: result,
	}
	err = json.NewDecoder(response.Body).Decode(&fullResp)
	if err != nil {
		return newUnableDecodedError(response.StatusCode, err)
	}
	if fullResp.Error != "" {
		return newServerError(fullResp.Error)
	}
	return nil
}

type fullRequest struct {
	Action  string `json:"action"`
	Version int    `json:"version"`
	Key     string `json:"key,omitempty"`
	Params  any    `json:"params,omitempty"`
}

type errResponse struct {
	Error string `json:"error"`
}

type fullResponse struct {
	Result any    `json:"result"`
	Error  string `json:"error"`
}
