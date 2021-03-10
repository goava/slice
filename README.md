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

## Lifecycle

- Initialize slice variables and components:
  - `slice.Info`: Application information: name, env and debug flag
  - `slice.Context`: Mutable context
  - `slice.Logger`: System logger (default: `stdout`)
  - `slice.ParameterParser`: Application parameter parser (default: `envconfig`)
- Create the container with user and slice components
- Parse all parameters (with bundle parameters)
- Invoke `BeforeStart` bundles hook
- Run dispatcher
- Invoke `BeforeShutdown` bundles hook

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
		slice.WithName("grpc-service"),
		slice.WithBundles(
			logging.Bundle,
			monitoring.Bundle,
			grpc.Bundle,
		),
		slice.WithComponents(
			di.Provide(NewDispatcher, di.As(new(slice.Dispatcher))),
			di.Provide(grpcsrv.NewService, di.As(new(grpc.Service))),
		),
	)
}
```

## Configuration

##### NAME


The name of application. Use `slice.WithName("your name")` to specify the
application name.

```go
slice.Run(
    slice.WithName("sliced"),
    // ...
)
```

##### ENV

The application environment. Use environment variable `ENV` to specify
application environment. The value can be any string.

## TODO

- [X] Environment bundle configuration
- [x] Configuration abstraction
- [ ] Batch update
- [ ] Replace `envconfig` to `parameter`
- [ ] Another parameter source (file, vault, consul).
- [ ] 90+% test coverage
- [ ] Bundle registry

## References

- [interface-based programming](https://en.wikipedia.org/wiki/Interface-based_programming)
- [modular programming](https://en.wikipedia.org/wiki/Modular_programming)
- [uber-go/fx](https://github.com/uber-go/fx)

