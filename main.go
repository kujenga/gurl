package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	v = flag.Bool("v", false, "enable verbose logging")
	m = flag.String("m", "GET", "HTTP method")

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

	resp, err := client().Do(req)
	check("error performing request:", err)

	fmt.Println("response:", resp)
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

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: t,
	}
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
	fmt.Printf("DEBUG: "+format, a...)
}
