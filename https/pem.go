package https

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func PemInit() {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization: []string{"Google Trust Services"},
		CommonName:   "GTS CA 1O1",
		Country:      []string{"US"},
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	path, _ := os.Getwd()
	//path = path + "/.pem/"
	//path =  "/opt/scn/.pem/"
	path = ".pem/"
	_ = os.Mkdir(path, 0777)
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	certOut, _ := os.Create(path + ".cert.pem")
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("failed to write data to cert.pem: %s", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("error closing cert.pem: %s", err)
	}

	keyOut, _ := os.Create(path + ".key.pem")
	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	_ = keyOut.Close()
}
