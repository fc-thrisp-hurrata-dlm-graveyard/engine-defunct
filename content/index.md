+++
title = "engine-index"
+++

Engine is a core package to drive a Go web framework. Routing, context, statuses
and more close the distance between the Go standard library and your own web
framework. 


### Install

    go get github.com/thrisp/engine

[*Quick Start*](/engine/quick)


### Documentaation

- [Current](/engine/documentation)

- [GoDoc](https://godoc.org/github.com/thrisp/engine)


### Configuration

Several configuration options are currently available, by specifying:


    engine.Option = value 


where needed in your code.

| Option | Explanation | Default |
| :---: | :---: | :---: |
| ServePanic | Serves a html page on panic | true |
| RedirectTrailingSlash | Enables automatic redirection if the current route can't be matched but a handler for the path with (without) the trailing slash exists | true |
| RedirectFixedPath | If enabled, the router tries to fix the current request path, if no handle is registered for it | true |
| HTMLStatus | All statuses send a simple html page | false |
| LoggingOn | All signals are sent to stdout through the logger or a default logger | false |
| MaxFormMemory | maximum size for file uploads, in bytes | 1000000 |