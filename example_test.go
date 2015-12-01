// Copyright 2015 The goprana Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goprana_test

import (
	"io/ioutil"
	"log"

	"github.com/jroimartin/goprana"
)

func ExampleClient_Get() {
	c := goprana.NewClient(goprana.DefaultPort, nil)

	resp, err := c.Get("gotquotes-go", "/api/quote/random")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s", output)
}
