# GoCertifi: SSL Certificates for Golang

This Go package contains a CA bundle that you can reference in your Go code.
This is useful for systems that do not have CA bundles that Golang can find
itself, or where a uniform set of CAs is valuable.

This is the same CA bundle that ships with the
[Python Requests](https://github.com/kennethreitz/requests) library, and is a
Golang specific port of [certifi](https://github.com/kennethreitz/certifi). The
CA bundle is derived from Mozilla's canonical set.

## Usage

You can use the gocertifi package as follows:

    import "github.com/certifi/gocertifi"
    cert_pool := gocertifi.CACerts()
