package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
	"github.com/sirupsen/logrus"
)

var toxiClient *toxiproxy.Client
var proxies []*toxiproxy.Proxy

// The following solution doesn't seem to work
// https://github.com/Shopify/toxiproxy/issues/175#issuecomment-301464691
// ./toxiproxy-cli create -l :26379 -u www.google.com:443 google
// curl -v -k -H "Host: www.google.com" https://127.0.0.1:26379
func main() {
	var err error
	toxiClient = toxiproxy.NewClient("127.0.0.1:8474")
	proxies, err = toxiClient.Populate([]toxiproxy.Proxy{{
		Name:     "google",
		Listen:   ":26379",
		Upstream: "www.google.com:443",
		Enabled: true,
	}})
	if err != nil {
		logrus.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodGet, "https://127.0.0.1:26379", nil)
	req.Host = "www.google.com"

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatal(err)
	}

	by, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(by))
}
