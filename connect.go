package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
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

	host := "l2e-forum.local"

	ips, err := net.LookupIP(host)
	if err == nil && len(ips) > 0 {
		host = ips[0].String()
	} else {
		host = "10.2.0.213"
	}

	config := &tls.Config{
		RootCAs:    rootCAs,
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", host+":443", config)

	if err != nil {
		log.Fatalf("Unable to establish secure connection!: %v", err)
	}

	return conn
}
