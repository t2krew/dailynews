package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:          800,
				MaxIdleConnsPerHost:   200,
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
				IdleConnTimeout:       90 * time.Second,
			},
		},
	}
}

func (c *Client) GetToQuery(url string, timeout time.Duration, headers map[string]string) (result map[string]string, err error) {
	result = make(map[string]string)
	_, ret, err := c.Get(url, timeout, headers)
	if err != nil {
		return result, err
	}

	result, err = QueryStringToMap(ret)

	if err != nil {
		return result, fmt.Errorf("false result format,decode failed, err %s, ret: %s", err.Error(), ret)
	}

	return result, nil
}

func (c *Client) GetToJSON(url string, timeout time.Duration, headers map[string]string) (result map[string]string, err error) {
	result = make(map[string]string)
	_, ret, err := c.Get(url, timeout, headers)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(ret), &result)

	if err != nil {
		return result, fmt.Errorf("false result format,decode failed, err %s, ret: %s", err.Error(), ret)
	}

	return result, nil
}

func (c *Client) Get(url string, timeout time.Duration, headers map[string]string) (int, string, error) {
	if c.client == nil {
		return 0, "", fmt.Errorf("%s", "httpClient not init")
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", fmt.Errorf("error while build request:  %v", err)
	}
	for k, v := range headers {
		if k == "Host" {
			req.Host = v
			continue
		}
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()

	req = req.WithContext(ctx)
	// try to get the response with url request
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	httpCode := resp.StatusCode

	// read the response body
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return httpCode, "", err
	}

	return httpCode, string(result), nil

}

func (c *Client) Post(url string, params []byte, timeout time.Duration, headers map[string]string) ([]byte, error) {
	if c.client == nil {
		return nil, fmt.Errorf("%s", "httpClient not init")
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for k, v := range headers {
		//stdlog.Print(k, v)
		request.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(request.Context(), timeout)
	defer cancel()

	request = request.WithContext(ctx)

	resp, err := c.client.Do(request)
	//stdlog.Printf("HttpPost response:%v", resp)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// read the response body
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
