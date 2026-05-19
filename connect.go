package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

func Connect() *tls.Conn {
	cert, err := os.ReadFile("cert.pem")
	if err != nil {
		log.Fatalf("Certificate Error: %s", err)
	}

	rootCAs := x509.NewCertPool()

	if ok := rootCAs.AppendCertsFromPEM(cert); !ok {
		log.Fatal("Failed to append Certificate to pool!")
	}

	config := &tls.Config{
		RootCAs:    rootCAs,
		ServerName: "localhost",
	}

	conn, err := tls.Dial("tcp", "12:443", config)

	if err != nil {
		log.Fatalf("Unable to establish secure connection!: %v", err)
	}

	return conn
}
