package imirhil

import (
    "net/url"
    "crypto/tls"
    "net/http"
    "log"
)

var (
    proxyURL *url.URL
)

func getProxy(req *http.Request) (uri *url.URL, err error) {
    uri, err = http.ProxyFromEnvironment(req)
    if err != nil {
        log.Printf("no proxy in environment")
        uri = &url.URL{}
    } else if uri == nil {
        log.Println("No proxy configured or url excluded")
    }
    return
}

func setupTransport(str string) (*http.Request, *http.Transport) {

    /*
       Proxy code taken from https://github.com/LeoCBS/poc-proxy-https/blob/master/main.go
    */
    myurl, err := url.Parse(str)
    if err != nil {
        log.Printf("error parsing %s: %v", str, err)
        return nil, nil
    }

    req, err := http.NewRequest("HEAD", str, nil)
    if err != nil {
        log.Printf("error: req is nil: %v", err)
        return nil, nil
    }
    req.Header.Set("Host", myurl.Host)
    req.Header.Add("User-Agent", "erc-checktls/proxyauth")

    // Get proxy URL
    proxyURL, err = getProxy(req)
    if ctx.proxyauth != "" {
        req.Header.Add("Proxy-Authorization", ctx.proxyauth)
    }

    transport := &http.Transport{
        Proxy:              http.ProxyURL(proxyURL),
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        ProxyConnectHeader: req.Header,
    }

    return req, transport
}


