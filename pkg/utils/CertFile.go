package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"path/filepath"

	"log"
	"math/big"
	"os"
	"time"
)

func CertFile() {
	// generate a new RSA key pair
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// create a template for the certificate
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"test"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"localhost"},
	}

	// create a self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatal(err)
	}

	// encode the certificate to PEM format
	certOut := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// encode the private key to PEM format
	keyOut := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	quicsDir := filepath.Join(homeDir, "quics")
	// write the certificate and the key to temporary files
	certFile, err := os.Create(filepath.Join(quicsDir, "cert-quics-cli.pem"))
	if err != nil {
		log.Fatal(err)
	}
	defer certFile.Close()

	keyFile, err := os.Create(filepath.Join(quicsDir, "key-quics-cli.pem"))
	if err != nil {
		log.Fatal(err)
	}
	defer keyFile.Close()

	if _, err := certFile.Write(certOut); err != nil {
		log.Fatal(err)
	}

	if _, err := keyFile.Write(keyOut); err != nil {
		log.Fatal(err)
	}

	// load the certificate and the key from the files
	// cert, err := tls.LoadX509KeyPair(certFile.Name(), keyFile.Name())
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
