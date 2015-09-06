package goprana_test

import (
	"io/ioutil"
	"log"

	"github.com/jroimartin/goprana"
)

func ExampleClient_Get() {
	c := goprana.NewClient(goprana.DefaultPort)

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
