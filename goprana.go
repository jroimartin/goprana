// Copyright 2015 The goprana Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package goprana is a Prana client for Go.
*/
package goprana

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// A Client represents a Prana client.
type Client struct {
	url     string
	httpcli *http.Client
}

// DefaultPort is default Prana port.
const DefaultPort = 8078

// NewClient returns Prana Client.
func NewClient(port int) Client {
	c := Client{
		url: fmt.Sprintf("http://localhost:%v", port),
		httpcli: &http.Client{
			Transport: &http.Transport{},
		},
	}
	return c
}

// DynamicProperties returns a map with the requested IDs and their values.
func (c Client) DynamicProperties(ids ...string) (props map[string]string, err error) {
	if len(ids) == 0 {
		return nil, errors.New("invalid ids value")
	}

	idsStr := ""
	for _, id := range ids {
		idsStr += fmt.Sprintf("id=%v&", id)
	}
	idsStr = idsStr[:len(idsStr)-1]

	url := fmt.Sprintf("%v/dynamicproperties?%v", c.url, idsStr)
	resp, err := c.httpcli.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is not 200 OK")
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&props)
	return
}

// Get sends a GET request to the specified VIP and path via Ribbon.
func (c Client) Get(vip, path string) (resp *http.Response, err error) {
	if vip == "" || path == "" {
		return nil, errors.New("invalid vip or path values")
	}

	url := fmt.Sprintf("%v/proxy?vip=%v&path=%v", c.url, vip, path)
	return c.httpcli.Get(url)
}

// Post sends a POST request to the specified VIP and path via Ribbon. Caller
// should close resp.Body when done reading from it. If the provided body is an
// io.Closer, it is closed after the request.
func (c Client) Post(vip, path, bodyType string, body io.Reader) (resp *http.Response, err error) {
	if vip == "" || path == "" {
		return nil, errors.New("invalid vip or path values")
	}

	url := fmt.Sprintf("%v/proxy?vip=%v&path=%v", c.url, vip, path)
	return c.httpcli.Post(url, bodyType, body)
}

// Do sends a request to the specified VIP and path via Ribbon. Caller should
// close resp.Body when done reading from it. If the provided body is an
// io.Closer, it is closed after the request.
func (c Client) Do(vip, path string, req *http.Request) (resp *http.Response, err error) {
	if vip == "" || path == "" {
		return nil, errors.New("invalid vip or path values")
	}

	rawurl := fmt.Sprintf("%v/proxy?vip=%v&path=%v", c.url, vip, path)
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req.URL = url
	return c.httpcli.Do(req)
}

// Hosts returns a list of hosts which are marked as UP in Eureka for the
// specific application and VIP passed as query parameters. If a VIP name is
// passed we filter the hosts matching that VIP name.
func (c Client) Hosts(appName, vip string) (hosts []string, err error) {
	if appName == "" {
		return nil, errors.New("invalid appName value")
	}

	url := fmt.Sprintf("%v/eureka/hosts?appName=%v", c.url, appName)
	if vip != "" {
		url += fmt.Sprintf("&vip=%v", vip)
	}
	resp, err := c.httpcli.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is not 200 OK")
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&hosts)
	return
}
