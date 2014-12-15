+++
title = "configuration"
+++
### Configuration

engine.New takes any number of configuration functions with the signature:

    func(\*Engine) error

These functions take the created engine to apply any configuration options, and returning any errors.

| Function | Explanation | Default |
| :---: | :---: | :---: |
| ServePanic(bool) | Serves a html page on panic | true |
| RedirectTrailingSlash(bool) | Enables automatic redirection if the current route can't be matched but a handler for the path with (without) the trailing slash exists | true |
| RedirectFixedPath(bool) | If enabled, the router tries to fix the current request path, if no handle is registered for it | true |
| HTMLStatus(bool) | All statuses send a simple html page | false |
| LoggingOn(bool) | All signals are sent to stdout through the logger or a default logger | false |
| Logger(\*log.Logger) | Sets logging on to true using the provided logger | nil |
| MaxFormMemory(int64) | maximum size for file uploads, in bytes | 1000000 |
