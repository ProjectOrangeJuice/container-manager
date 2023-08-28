package connection

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Println("Failed to cast connection to TLS connection")
		return
	}
	// Complete the handshake to get the client certificate
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake error: %v", err)
		return
	}

	certs := tlsConn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		log.Println("Client presented no certificates.")
		return
	}
	clientCert := certs[0]

	if !isCertificatePinned(clientCert) {
		if !promptAllowConnection(clientCert) {
			log.Println("Connection denied.")
			return
		}
	}

	log.Println("Connection allowed.")
}

var pinnedCertificates = map[string]bool{}

func isCertificatePinned(cert *x509.Certificate) bool {
	return pinnedCertificates[cert.SerialNumber.String()]
}

func promptAllowConnection(cert *x509.Certificate) bool {
	fmt.Printf("Unknown client certificate: %s\n", cert.Subject.CommonName)
	fmt.Print("Allow connection? [Y/n]: ")
	var response string
	_, _ = fmt.Scanln(&response)
	return response == "Y" || response == "y"
}
