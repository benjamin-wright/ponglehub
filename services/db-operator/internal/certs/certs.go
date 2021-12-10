package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

func GenerateCACerts() ([]byte, []byte, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"PongleHub, INC."},
			Country:       []string{"GB"},
			Province:      []string{""},
			Locality:      []string{"Pongletown"},
			StreetAddress: []string{"Snowhere"},
			PostalCode:    []string{"99999"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate CA key: %+v", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CA certificate: %+v", err)
	}

	certPEM := new(bytes.Buffer)
	if err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return nil, nil, fmt.Errorf("failed to PEM encode CA certificate: %+v", err)
	}

	keyPEM := new(bytes.Buffer)
	if err = pem.Encode(keyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return nil, nil, fmt.Errorf("failed to PEM encode CA key: %+v", err)
	}

	return certPEM.Bytes(), keyPEM.Bytes(), nil
}

func GenerateNodeCerts(dns []string, cacert []byte, cakey []byte) ([]byte, []byte, error) {
	block, _ := pem.Decode(cacert)
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing cacert from byte data: %+v", err)
	}

	block, _ = pem.Decode(cakey)
	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing cakey from byte data: %+v", err)
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			Organization:  []string{"PongleHub, INC."},
			Country:       []string{"GB"},
			Province:      []string{""},
			Locality:      []string{"Pongletown"},
			StreetAddress: []string{"Snowhere"},
			PostalCode:    []string{"99999"},
			CommonName:    "node",
		},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
		DNSNames:    dns,
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate cert key: %+v", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &key.PublicKey, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate certificate: %+v", err)
	}

	certPEM := new(bytes.Buffer)
	if err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return nil, nil, fmt.Errorf("failed to PEM encode certificate: %+v", err)
	}

	keyPEM := new(bytes.Buffer)
	if err = pem.Encode(keyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return nil, nil, fmt.Errorf("failed to PEM encode key: %+v", err)
	}

	return certPEM.Bytes(), keyPEM.Bytes(), nil
}
