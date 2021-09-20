/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package gocertifi

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"
)

func TestGetCerts(t *testing.T) {
	certPool, err := CACerts()
	if certPool == nil || err != nil || len(certPool.Subjects()) == 0 {
		t.Errorf("Failed to return the certificates.")
	}
}

func parsePEM(pemCerts []byte) (certs []*x509.Certificate, err error) {
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	return
}

func checkRootCertsPEM(t *testing.T, pemCerts []byte, when time.Time) (ok bool) {
	now := time.Now()
	t.Logf("Checking certificate validity on %s...", when)
	certs, err := parsePEM(pemCerts)
	if err != nil {
		t.Error(err)
		return false
	}

	roots := x509.NewCertPool()
	for _, cert := range certs {
		roots.AddCert(cert)
	}

	var minExpires time.Time
	ok = true
	for _, cert := range certs {
		if !cert.IsCA {
			t.Errorf("\u274C %s: not a certificate authority", cert.Subject)
		}
		// This check of keyusage is based on my understanding of key usage purposes
		// https://cabforum.org/wp-content/uploads/CA-Browser-Forum-BR-1.8.0.pdf
		// Section 1.4.2 has no usage restrictions
		if cert.KeyUsage&(x509.KeyUsage(-1)^(x509.KeyUsageCertSign|x509.KeyUsageCRLSign|x509.KeyUsageDigitalSignature)) != 0 {
			t.Logf("\u26A0 %s key usage %#x (see constants at https://golang.org/pkg/crypto/x509/#KeyUsage)", cert.Subject, cert.KeyUsage)
		} else if cert.KeyUsage&(x509.KeyUsageCertSign|x509.KeyUsageCRLSign) == 0 {
			// If the certificate authority is not allowed to sign certificates, why is it here?
			// https://cabforum.org/baseline-requirements-certificate-contents/#CA-Certificates
			t.Logf("\u26A0 %s key usage %#x (see constants at https://golang.org/pkg/crypto/x509/#KeyUsage)", cert.Subject, cert.KeyUsage)
		}
		if minExpires.IsZero() || cert.NotAfter.Before(minExpires) {
			minExpires = cert.NotAfter
		}
		// Check that the certificate is valid now
		if cert.NotBefore.After(now) {
			t.Errorf("\u274C %s: fails NotBefore check: %s", cert.Subject, cert.NotBefore)
			continue
		}
		if cert.NotAfter.Before(now) {
			t.Errorf("\u274C %s: fails NotAfter check: \033[31m%s\033[m", cert.Subject, cert.NotAfter)
			// ... and that it will still be valid later
		} else if cert.NotAfter.Before(when) {
			t.Logf("\u26A0 %s: fails NotAfter check: \033[31m%s\033[m", cert.Subject, cert.NotAfter)
			continue
		}
		_, err := cert.Verify(x509.VerifyOptions{
			Roots:       roots,
			CurrentTime: when,
		})
		if err != nil {
			t.Errorf("\u274C %s: %s", cert.Subject, err)
			ok = false
		} else {
			t.Logf("\u2705 %s (expires: %s)", cert.Subject, cert.NotAfter)
		}
	}
	if ok {
		t.Log("Success.")
		t.Logf("MinExpire: %s", minExpires)
	}
	return
}

func TestCerts(t *testing.T) {
	// Check that certificates will still be valid in 3 months
	checkRootCertsPEM(t, []byte(pemcerts), time.Now().AddDate(0, 3, 0))
}
