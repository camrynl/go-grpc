package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"

	"github.com/cilium/cilium/api/v1/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	hubblePort     = 4244
	clientCertFile = "./tls.crt"
	clientKeyFile  = "./tls.key"
)

func main() {
	serverAddr := fmt.Sprintf("localhost:%d", hubblePort)
	tlsConfig := tls.Config{
		InsecureSkipVerify: false,
		ServerName:         serverAddr,
	}
	var cert *tls.Certificate
	if clientCertFile != "" && clientKeyFile != "" {
		c, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			log.Fatalf("failed to load keypair: %v", err)
		}
		cert = &c
	}
	tlsConfig.GetClientCertificate = func(_ *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		if cert == nil {
			log.Printf("certs not found")
			return nil, errors.New("mTLS client certificate requested, but not provided")
		}
		return cert, nil
	}

	creds := credentials.NewTLS(&tlsConfig)

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := observer.NewObserverClient(conn)
	flows, err := client.GetFlows(context.Background(), &observer.GetFlowsRequest{})
	if err != nil {
		log.Fatalf("failed to get flows: %v", err)
	}
	log.Printf("flows: %v", flows)
}
