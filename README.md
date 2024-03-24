# Fillmore Labs Async

[![Go Reference](https://pkg.go.dev/badge/fillmore-labs.com/async.svg)](https://pkg.go.dev/fillmore-labs.com/async)
[![Build status](https://badge.buildkite.com/88d2f145eee0fde273b7bdbe9e95ee7eeee6e0e48b443ff227.svg)](https://buildkite.com/fillmore-labs/async)
[![GitHub Workflow](https://github.com/fillmore-labs/async/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/fillmore-labs/async/actions/workflows/test.yml)
[![Test Coverage](https://codecov.io/gh/fillmore-labs/async/graph/badge.svg?token=B70VNID5KK)](https://codecov.io/gh/fillmore-labs/async)
[![Maintainability](https://api.codeclimate.com/v1/badges/edf5df13e2ef438663af/maintainability)](https://codeclimate.com/github/fillmore-labs/async/maintainability)
[![Go Report Card](https://goreportcard.com/badge/fillmore-labs.com/async)](https://goreportcard.com/report/fillmore-labs.com/async)
[![License](https://img.shields.io/github/license/fillmore-labs/async)](https://www.apache.org/licenses/LICENSE-2.0)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ffillmore-labs%2Fasync.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Ffillmore-labs%2Fasync)

The `async` package provides interfaces and utilities for writing asynchronous code in Go.

## Motivation

...

## Usage

Assuming you have a synchronous function `func getMyIP(ctx context.Context) (string, error)` returning your external IP
address (see [GetMyIP](#getmyip) for an example).

Now you can do

```go
package main

import (
	"context"
	"log/slog"
	"time"

	"fillmore-labs.com/async"
)

func main() {
	const timeout = 2 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	future := async.NewAsync(func() (string, error) { return getMyIP(ctx) })

	// other queries

	if ip, err := future.Await(ctx); err == nil {
		slog.Info("Found IP", "ip", ip)
	} else {
		slog.Error("Failed to fetch IP", "error", err)
	}
}
```

decoupling query construction from result processing.

### GetMyIP

Sample code to retrieve your IP address:

```go
package main

import (
	"context"
	"encoding/json"
	"net/http"
)

const serverURL = "https://httpbin.org/ip"

func getMyIP(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	ipResponse := struct {
		Origin string `json:"origin"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&ipResponse)

	return ipResponse.Origin, err
}
```

## Links

- [Futures and Promises](https://en.wikipedia.org/wiki/Futures_and_promises) in the English Wikipedia
