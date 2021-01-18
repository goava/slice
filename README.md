Slice (Work in progress)
========================

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://pkg.go.dev/github.com/goava/slice)
[![Release](https://img.shields.io/github/tag/goava/slice.svg?label=release&color=24B898&logo=github&style=for-the-badge)](https://github.com/goava/slice/releases/latest)
[![Build Status](https://img.shields.io/travis/goava/slice.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/goava/slice)
[![Code Coverage](https://img.shields.io/codecov/c/github/goava/slice.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/goava/slice)

## Problem

During the process of writing software in the team, you develop a
certain style and define standards that meet the requirements for this
software. These standards grow into libraries and frameworks. This is our
approach based on
[interface-based programming](https://en.wikipedia.org/wiki/Interface-based_programming)
and
[modular programming](https://en.wikipedia.org/wiki/Modular_programming).

## Overview

```go
package main

import (
	"github.com/goava/slice"
	"github.com/goava/di"
	// imports omitted
)

func main() {
	slice.Run(
		slice.SetName("grpc-service"),
		slice.RegisterBundles(
			logging.Bundle(),
			monitoring.Bundle(),
			grpc.Bundle(),
		),
		slice.ConfigureContainer(
			di.Provide(NewDispatcher, di.As(new(slice.Dispatcher))),
			di.Provide(grpcsrv.NewService, di.As(new(grpc.Service))),
		),
	)
}
```

## Configuration

##### NAME


The name of application. Use `slice.SetName("your name")` to specify the
application name.

```go
slice.Run(
    slice.SetName("sliced"),
    // ...
)
```

##### ENV

The application environment. Use environment variable `ENV` to specify
application environment. The value can be any string.

## TODO

- [X] Environment bundle configuration
- [ ] Configuration abstraction
- [ ] 90+% test coverage
- [ ] Bundle registry
- [ ] Registry interface

## References

- [interface-based programming](https://en.wikipedia.org/wiki/Interface-based_programming)
- [modular programming](https://en.wikipedia.org/wiki/Modular_programming)
- [uber-go/fx](https://github.com/uber-go/fx)

