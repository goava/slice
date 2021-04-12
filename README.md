Slice (Work in progress)
========================

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://pkg.go.dev/github.com/goava/slice)
[![Release](https://img.shields.io/github/tag/goava/slice.svg?label=release&color=24B898&logo=github&style=for-the-badge)](https://github.com/goava/slice/releases/latest)
[![Build Status](https://img.shields.io/travis/goava/slice.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/goava/slice)
[![Code Coverage](https://img.shields.io/codecov/c/github/goava/slice.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/goava/slice)

## Problem

During the process of writing software in the team, you develop a
certain style and define standards that meet the requirements for this
software. These standards grow into libraries and frameworks. This is
our approach based on
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

## Minimal start

The minimum that you need to run the application on `slice`.

### Set application name

Use `slice.WithName("your name")` to specify the application name.

```go
slice.Run(
    slice.WithName("sliced"),
    // ...
)
```

### Check application environment

Use environment variable `ENV` to specify the application environment.
The value can be any string. Environments that have a prefix `dev` will
be recognized as a development environment. Others will be recognized as
production.

### Provide application dispatcher

Provide your own `slice.Dispatcher` implementation:

```go
slice.Run(
    slice.WithName("sliced"),
    slice.WithComponents(
        di.Provide(NewDispatcher, di.As(new(slice.Dispatcher))),
    ),
)
```

## Lifecycle

- Initialize slice variables and components:
  - `slice.Info`: Application information: name, env and debug flag
  - `slice.Context`: Application context
  - `slice.Logger`: System logger (default: `stdout`)
  - `slice.ParameterParser`: Application parameter parser (default:
    `envconfig`)
- Create the container with user and slice components
- Parse all parameters (with bundle parameters)
- Invoke `BeforeStart` bundles hook
- Run dispatcher
- Invoke `BeforeShutdown` bundles hook

## Lifecycle details

### Parameter parsing

Applications created with `slice` support parameter parsing. By default,
it's processed by
[envconfig](https://github.com/kelseyhightower/envconfig).

You can use your own parameter parser. To do this, implement the
`ParameterParser` interface and use it using the `WithParameterParser()`
or by setting the `ParameterParser` field of the `Application`.

You can print all parameters by using `<binary-name> --parameters`.
Example output with default parameter parser and following structure:

```go
// Parameters contains application configuration.
type Parameters struct {
	Addr         string        `envconfig:"addr" required:"true" desc:"Server address"`
	ReadTimeout  time.Duration `envconfig:"read_timeout" required:"true" desc:"Server read timeout"`
	WriteTimeout time.Duration `envconfig:"write_timeout" required:"true" desc:"Server write timeout"`
}
```

Output:

```text
KEY              TYPE        DEFAULT    REQUIRED    DESCRIPTION
ADDR             String                 true        Server address
READ_TIMEOUT     Duration               true        Server read timeout
WRITE_TIMEOUT    Duration               true        Server write timeout
```

## References

- [interface-based programming](https://en.wikipedia.org/wiki/Interface-based_programming)
- [modular programming](https://en.wikipedia.org/wiki/Modular_programming)
- [uber-go/fx](https://github.com/uber-go/fx)

