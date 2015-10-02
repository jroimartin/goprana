// Copyright 2015 The goprana Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goprana_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jroimartin/goprana"
)

const (
	appName = "go-gotquotes"
	vip     = "go-gotquotes"
)

type quoteResp struct {
	Quote   string
	Counter int
}

var reqChecks = map[string]quoteResp{
	"/api/quote/1": quoteResp{
		Quote: "You know nothing, Jon Snow.",
	},
}

var propChecks = map[string]string{
	"eureka.vipAddress": "go-gotquotes",
	"eureka.port":       "5000",
}

func TestDynamicProperties(t *testing.T) {
	c := goprana.NewClient(goprana.DefaultPort)

	ids := []string{}
	for k := range propChecks {
		ids = append(ids, k)
	}

	props, err := c.DynamicProperties(ids...)
	if err != nil {
		t.Fatal(err)
	}

	for id, want := range propChecks {
		val, ok := props[id]
		if !ok {
			t.Errorf("property %v not found", id)
			continue
		}
		if val != want {
			t.Errorf("got property %v=%v, want %v", id, val, want)
		}
	}
}

func TestGet(t *testing.T) {
	c := goprana.NewClient(goprana.DefaultPort)

	for path, want := range reqChecks {
		resp, err := c.Get(vip, path)
		if err != nil {
			t.Fatal(err)
		}

		dec := json.NewDecoder(resp.Body)
		var r quoteResp
		if err := dec.Decode(&r); err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if r.Quote != want.Quote {
			t.Errorf("got quote %v, want %v", r.Quote, want.Quote)
		}
		if r.Counter <= 0 {
			t.Errorf("got quote counter %v, want >0", r.Counter)
		}
	}
}

func TestPost(t *testing.T) {
	c := goprana.NewClient(goprana.DefaultPort)

	for path := range reqChecks {
		resp, err := c.Post(vip, path, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusNotImplemented {
			t.Errorf("got status %v, want 501 Not Implemented", resp.StatusCode)
		}
	}
}

func TestDo(t *testing.T) {
	c := goprana.NewClient(goprana.DefaultPort)

	for path, want := range reqChecks {
		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := c.Do(vip, path, req)
		if err != nil {
			t.Fatal(err)
		}

		dec := json.NewDecoder(resp.Body)
		var r quoteResp
		if err := dec.Decode(&r); err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if r.Quote != want.Quote {
			t.Errorf("got quote %v, want %v", r.Quote, want.Quote)
		}
		if r.Counter <= 0 {
			t.Errorf("got quote counter %v, want >0", r.Counter)
		}
	}
}

func TestHosts(t *testing.T) {
	c := goprana.NewClient(goprana.DefaultPort)

	hosts, err := c.Hosts(appName, vip)
	if err != nil {
		t.Fatal(err)
	}

	if len(hosts) == 0 {
		t.Error("got 0 hosts")
	}
}
