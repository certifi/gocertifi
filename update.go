package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	// Retrieve cacerts.pem
	resp, err := http.Get("http://ci.kennethreitz.org/job/ca-bundle/lastSuccessfulBuild/artifact/cacerts.pem")
	if err != nil {
		panic("retrieval of cacerts.pem failed")
	}
	defer resp.Body.Close()

	// Read the cacerts.pem file into memory
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("reading of response body failed")
	}

	// Open our certifi.go file
	file, err := ioutil.ReadFile("certifi.go")
	if err != nil {
		panic("error reading certifi.go")
	}

	// Split certifi.go on the string containing the certificates
	certFile := strings.Split(string(file), "`\n")
	if len(certFile) != 3 {
		panic("error splitting certifi.go")
	}
	// Replace the old certs string with the contents of cacerts.pem
	certFile[1] = string(body)
	// Join the file parts into a single []byte, for writing to disk
	outBytes := []byte(strings.Join(certFile, "`\n"))

	// Overwrite certifi.go with our updated contents
	err = ioutil.WriteFile("certifi.go", outBytes, 0644)
	if err != nil {
		panic("error writing to certifi.go")
	}
}
