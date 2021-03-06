// +build ignore

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
	"bufio"
	"bytes"
)

func main() {
	if len(os.Args) != 2 || !strings.HasPrefix(os.Args[1], "https://") {
		log.Fatal("usage: go run gen.go <url>")
	}
	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("expected 200, got", resp.StatusCode)
	}
	defer resp.Body.Close()

	var bundle bytes.Buffer

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) == 0 || b[0] == '#' {
			continue
		}
		bundle.Write(b)
		bundle.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("failed to read response body fully", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(bundle.Bytes()) {
		log.Fatalf("can't parse cerficiates from %s", url)
	}

	fp, err := os.Create("certifi.go")
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	tmpl.Execute(fp, struct {
		Timestamp time.Time
		URL       string
		Bundle    string
	}{
		Timestamp: time.Now(),
		URL:       url,
		Bundle:    bundle.String(),
	})
}

var tmpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// {{ .Timestamp }}
// {{ .URL }}

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package gocertifi

//go:generate go run gen.go "{{ .URL }}"

import "crypto/x509"

const pemcerts string = ` + "`" + `
{{ .Bundle }}
` + "`" + `

// CACerts builds an X.509 certificate pool containing the
// certificate bundle from {{ .URL }} fetch on {{ .Timestamp }}.
// Will never actually return an error.
func CACerts() (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(pemcerts))
	return pool, nil
}
`))
