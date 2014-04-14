package gocertifi

import "testing"

func TestGetCerts(t *testing.T) {
	cert_pool := CACerts()
	if cert_pool == nil {
		t.Errorf("Failed to return the certificates.")
	}
}
