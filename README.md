# GoCertifi: SSL Certificates for Golang

This Go package contains a CA bundle that you can reference in your Go code.
This is useful for systems that do not have CA bundles that Golang can find
itself, or where a uniform set of CAs is valuable.

This is the same CA bundle that ships with the
[Python Requests](https://github.com/kennethreitz/requests) library, and is a
Golang specific port of [certifi](https://github.com/kennethreitz/certifi). The
CA bundle is derived from Mozilla's canonical set.

## Usage

You can use the `gocertifi` package as follows:

```go
import "github.com/certifi/gocertifi"
cert_pool, err := gocertifi.CACerts()
```

## Detailed Documentation

Import as follows:

```go
import "github.com/certifi/gocertifi"
```

### Errors

```go
var ErrParseFailed = errors.New("gocertifi: error when parsing certificates")
```

### Functions

```go
func CACerts() (*x509.CertPool, error)
```
CACerts builds an X.509 certificate pool containing the Mozilla CA Certificate
bundle. Returns nil on error along with an appropriate error code.
