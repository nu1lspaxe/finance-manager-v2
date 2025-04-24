package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var certFile, keyFile, caFile string

func init() {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	certFile = viper.GetString("certs.path.cert")
	keyFile = viper.GetString("certs.path.key")
	caFile = viper.GetString("certs.path.ca")

	// err := setupCertificate()
	// if err != nil {
	// 	panic(err)
	// }
}

func LoadTLSConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA cert")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"http/1.1", "h2"},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      caCertPool,
	}, nil
}

// func setupCertificate() error {
// 	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	if err != nil {
// 		return err
// 	}

// 	notBefore := time.Now()
// 	notAfter := notBefore.AddDate(1, 0, 0)

// 	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
// 	if err != nil {
// 		return err
// 	}

// 	template := x509.Certificate{
// 		SerialNumber: serialNumber,
// 		Subject: pkix.Name{
// 			Organization: []string{"Finance Manager"},
// 			CommonName:   "localhost",
// 		},
// 		DNSNames:  []string{"localhost"},
// 		NotBefore: notBefore,
// 		NotAfter:  notAfter,
// 		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
// 		ExtKeyUsage: []x509.ExtKeyUsage{
// 			x509.ExtKeyUsageServerAuth,
// 		},
// 		BasicConstraintsValid: true,
// 		IsCA:                  true,
// 		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
// 	}

// 	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &private.PublicKey, private)
// 	if err != nil {
// 		return fmt.Errorf("failed to create certificate: %v", err)
// 	}

// 	certOut, err := os.OpenFile(certFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		return fmt.Errorf("failed to open cert file: %v", err)
// 	}
// 	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert}); err != nil {
// 		certOut.Close()
// 		return fmt.Errorf("failed to encode cert: %v", err)
// 	}
// 	certOut.Close()

// 	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		return fmt.Errorf("failed to open key file: %v", err)
// 	}
// 	privateBytes, err := x509.MarshalECPrivateKey(private)
// 	if err != nil {
// 		keyOut.Close()
// 		return fmt.Errorf("failed to marshal private key: %v", err)
// 	}
// 	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privateBytes}); err != nil {
// 		keyOut.Close()
// 		return fmt.Errorf("failed to encode key: %v", err)
// 	}
// 	keyOut.Close()

// 	return nil
// }
