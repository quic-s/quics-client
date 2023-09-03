package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"path/filepath"

	"log"
	"math/big"
	"os"

	"github.com/quic-s/quics-client/pkg/viper"
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
	}

	// create a self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatal(err)
	}

	// encode the certificate, key to PEM format
	certOut := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyOut := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// write the certificate and key to disk
	quicsDir := GetQuicsDirPath()
	certFile, err := os.Create(filepath.Join(quicsDir, viper.GetViperEnvVariables("QUICS_CLI_CERT_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	defer certFile.Close()

	keyFile, err := os.Create(filepath.Join(quicsDir, viper.GetViperEnvVariables("QUICS_CLI_KEY_NAME")))
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
