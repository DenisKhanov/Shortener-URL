// Package https provides functionality for generating self-signed TLS certificates and private keys,
// and saving them to files for use in HTTPS server configurations.
package https

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/sirupsen/logrus"
	"math/big"
	"net"
	"os"
	"time"
)

// HTTPS represents the HTTPS configuration containing the generated certificate and private key.
type HTTPS struct {
	cert       []byte
	privateKey *rsa.PrivateKey
}

// NewHTTPS generates a self-signed TLS certificate and private key, saves them to files,
// and returns the HTTPS configuration.
func NewHTTPS() (*HTTPS, error) {

	var cert = &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"DKShortenerPJT"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logrus.WithError(err).Error("failed to generate private key")
		return nil, err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		logrus.WithError(err).Error("failed to create certificate")
		return nil, err
	}

	newHTTPS := &HTTPS{
		cert:       certBytes,
		privateKey: privateKey,
	}
	if err = newHTTPS.saveCertToFile(); err != nil {
		return nil, err
	}
	if err = newHTTPS.saveKeyToFile(); err != nil {
		return nil, err
	}
	return newHTTPS, nil
}

// saveCertToFile saves the generated certificate to a file.
func (h *HTTPS) saveCertToFile() error {
	var certPEM bytes.Buffer
	if err := pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: h.cert,
	}); err != nil {
		logrus.WithError(err).Errorf("filed to create %s file", models.CertPEM)
		return err
	}
	err := os.WriteFile(models.CertPEM, certPEM.Bytes(), 0644)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// saveKeyToFile saves the generated private key to a file.
func (h *HTTPS) saveKeyToFile() error {
	var privateKeyPEM bytes.Buffer
	if err := pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(h.privateKey),
	}); err != nil {
		logrus.WithError(err).Errorf("filed to create %s file", models.PrivateKeyPEM)
		return err
	}
	err := os.WriteFile(models.PrivateKeyPEM, privateKeyPEM.Bytes(), 0644)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
