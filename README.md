# clientip

[![GoDoc](https://godoc.org/github.com/dwisiswant0/clientip?status.svg)](http://godoc.org/github.com/dwisiswant0/clientip)
[![Go Report Card](https://goreportcard.com/badge/github.com/dwisiswant0/clientip)](https://goreportcard.com/report/github.com/dwisiswant0/clientip)

Go library to retrieve the client IP from HTTP requests.

> History:
> * This library is a fork of [victorkt/clientip](https://github.com/victorkt/clientip) with minimal refactor.
> * This library is a port of [pbojinov/request-ip](https://github.com/pbojinov/request-ip) with additional tooling for Go servers.

## Installation

```shell script
$ go get github.com/dwisiswant0/clientip
```

## Basic usage

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dwisiswant0/clientip"
)

func main() {
	http.HandleFunc("/", HelloServer)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	ip := clientip.FromRequest(r)
	fmt.Fprintf(w, "Hello, your IP is %s!", ip)
}

```

## Middleware

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dwisiswant0/clientip"
)

func main() {
	handler := http.HandlerFunc(HelloServer)
	http.Handle("/", clientip.Middleware(handler))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	ip := clientip.FromContext(r.Context())
	fmt.Fprintf(w, "Hello, your IP is %s!", ip)
}
```

## How it works

The library will try to get the client IP from a list of headers and falls back on request.RemoteAddr

The order in which the headers are checked is:

1. X-Client-IP
2. X-Forwarded-For (Header may return multiple IP addresses in the format: "client IP, proxy 1 IP, proxy 2 IP", so we take the the first one.)
3. CF-Connecting-IP (Cloudflare)
4. Fastly-Client-Ip (Fastly CDN and Firebase hosting header when forwared to a cloud function)
5. True-Client-Ip (Akamai and Cloudflare)
6. X-Real-IP (Nginx proxy/FastCGI)
7. X-Cluster-Client-IP (Rackspace LB, Riverbed Stingray)
8. X-Forwarded and Forwarded-For (Variations of #2)
9. request.RemoteAddr

If a valid IP was found and it contains a port number, the port will be ignored. If no valid IP is found, it returns a `nil` `net.IP`

# License

The MIT License (MIT) - 2020