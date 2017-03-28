package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	v      = flag.Bool("v", false, "enable verbose logging")
	m      = flag.String("m", "GET", "HTTP method")
	follow = flag.Bool("follow", false, "follow redirects")

	h = flag.Bool("h", false, "display help information")

	insecure = flag.Bool("insecure", false, "allow insecure HTTPS connections")
	timeout  = flag.String("timeout", "10s", "timeout for the http client")
)

func main() {
	flag.Parse()

	switch {
	case *h:
		flag.PrintDefaults()
		os.Exit(0)
	}

	urlStr := flag.Arg(0)

	if urlStr == "" {
		nope("url argument must be present")
	}

	req, err := http.NewRequest(*m, urlStr, nil)
	check("invalid url:", err)
	printf("> %+v\n", req)

	resp, err := client().Do(req)
	check("error performing request:", err)
	printf("< %+v\n", resp)

	io.Copy(os.Stdout, resp.Body)
}

func client() *http.Client {
	t, err := time.ParseDuration(*timeout)
	check("parsing timeout:", err)

	var tlsConfig *tls.Config
	if *insecure {
		debug("tls configured with InsecureSkipVerify")
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: t,
	}

	if !*follow {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}

func check(msg string, err error) {
	if err != nil {
		debugf("error (%T): %#v\n", err, err)
		nope(msg, err)
	}
}

func nope(a ...interface{}) {
	fmt.Println(a...)
	os.Exit(1)
}

func printf(format string, a ...interface{}) {
	if !*v {
		return
	}
	fmt.Printf(format, a...)
}

func debug(a ...interface{}) {
	if !*v {
		return
	}
	fmt.Println(append([]interface{}{"DEBUG:"}, a...)...)
}
func debugf(format string, a ...interface{}) {
	if !*v {
		return
	}
	printf("DEBUG: "+format, a...)
}
