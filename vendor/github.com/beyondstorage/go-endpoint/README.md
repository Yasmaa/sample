# go-endpoint

[![Build Status](https://github.com/beyondstorage/go-endpoint/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/beyondstorage/go-endpoint/actions?query=workflow%3A%22Unit+Test%22)
[![Go dev](https://pkg.go.dev/badge/github.com/beyondstorage/go-endpoint)](https://pkg.go.dev/github.com/beyondstorage/go-endpoint)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/beyondstorage/go-endpoint/blob/master/LICENSE)
[![go storage dev](https://img.shields.io/matrix/go-endpoint:aos.dev.svg?server_fqdn=chat.aos.dev&label=%23go-endpoint%3Aaos.dev&logo=matrix)](https://matrix.to/#/#go-endpoint:aos.dev) <!-- Need update after matrix updated -->

Both human and machine readable  endpoint format.

## Format

```
<protocol>:<value>+
```

For example:

- File: `file:/var/cache/data`
- HTTP: `http:example.com:80`
- HTTPS: `https:example.com:443`

## Quick Start

```go
ep, err := endpoint.Parse("https:example.com")
if err != nil {
	log.Fatal("parse: ", err)
}

switch ep.Protocol() {
case ProtocolHTTP:
    url, host, port := ep.HTTP()
    log.Println("url: ", url)
    log.Println("host: ", host)
    log.Println("port: ", port)
case ProtocolHTTPS:
    url, host, port := ep.HTTPS()
    log.Println("url: ", url)
    log.Println("host: ", host)
    log.Println("port: ", port)
case ProtocolFile:
    path := ep.File()
    log.Println("path: ", path)
default:
    panic("unsupported protocol")
}
```
