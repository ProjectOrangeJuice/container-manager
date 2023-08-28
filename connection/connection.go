package connection

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"

	"github.com/ProjectOrangeJuice/vm-manager-server/serverConfig"
)

func Setup() Clients {
	return &allClients{}
}

type Clients interface {
	HandleClient(conn net.Conn)
	GetAllClients() []*Client
}

func InitFingerprints(fingerprints []serverConfig.Fingerprint) {
	for _, fingerprint := range fingerprints {
		pinnedCertificates[fingerprint.Fingerprint] = fingerprint.AllowConnect
	}

}

func (ac *allClients) HandleClient(conn net.Conn) {
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
	ac.AddClient(clientCert.Subject.CommonName, clientCert.SerialNumber.String(), conn)
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
	if response == "Y" || response == "y" {
		addFingerPrint(cert.SerialNumber.String(), cert.SerialNumber.String(), cert.Subject.CommonName, true)
	}
	return response == "Y" || response == "y" // This will be on the UI, but fornow we accept an input through the terminal
}

func addFingerPrint(fingerprint, serial, name string, allow bool) {
	pinnedCertificates[fingerprint] = allow
	serverConfig.AddFingerprint(serverConfig.Fingerprint{
		Name:         name,
		SerialNumber: serial,
		Fingerprint:  fingerprint, // This is the serial number for now
		AllowConnect: allow,
	})
}
