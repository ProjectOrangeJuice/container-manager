package connection

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"time"

	"github.com/ProjectOrangeJuice/vm-manager-server/serverConfig"
)

var pinnedCertificates = map[string]bool{}

func Setup(config serverConfig.Config) Clients {
	// Load the accepted clients
	ac := &allClients{}
	ac.InitFingerprints(config.ClientFingerprints)
	return ac
}

type Clients interface {
	HandleClient(conn net.Conn)
	GetActiveClients() []*Client
	GetWaitingClients() []ClientDetails
	GetAcceptedClients() []ClientDetails
}

func (ac *allClients) InitFingerprints(fingerprints []serverConfig.Fingerprint) {
	for _, fingerprint := range fingerprints {
		pinnedCertificates[fingerprint.Fingerprint] = fingerprint.AllowConnect
		ac.AcceptedClients = append(ac.AcceptedClients, ClientDetails{
			Name:        fingerprint.Name,
			Serial:      fingerprint.SerialNumber,
			Fingerprint: fingerprint.Fingerprint,
		})
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
		if !ac.promptAllowConnection(clientCert) {
			log.Println("Connection denied.")
			return
		}
	}

	log.Println("Connection allowed.")
	ac.AddClient(clientCert.Subject.CommonName, clientCert.SerialNumber.String(), conn)
}

func (ac *allClients) promptAllowConnection(cert *x509.Certificate) bool {

	ac.clientLock.Lock()
	ac.WaitingClients = append(ac.WaitingClients, ClientDetails{
		Name:        cert.Subject.CommonName,
		Serial:      cert.SerialNumber.String(),
		Fingerprint: cert.SerialNumber.String(),
	})

	ac.clientLock.Unlock()

	// Every 10 seconds check list, after 3 hours, stop checking
	for i := 0; i < 1080; i++ {
		time.Sleep(10 * time.Second)

		// Check if the client is in waiting list
		ac.clientLock.Lock()
		found := false
		for _, client := range ac.WaitingClients {
			if client.Serial == cert.SerialNumber.String() {
				found = true
				break
			}
		}
		if !found {
			// check if the client is in accepted list
			for _, client := range ac.AcceptedClients {
				if client.Serial == cert.SerialNumber.String() {
					ac.clientLock.Unlock()
					ac.addFingerPrint(cert.SerialNumber.String(), cert.SerialNumber.String(), cert.Subject.CommonName, true)
					return true
				}
			}
			ac.clientLock.Unlock()
			ac.addFingerPrint(cert.SerialNumber.String(), cert.SerialNumber.String(), cert.Subject.CommonName, false)
			return false
		}
		ac.clientLock.Unlock()
	}
	return false
}

func (ac *allClients) addFingerPrint(fingerprint, serial, name string, allow bool) {
	pinnedCertificates[fingerprint] = allow
	serverConfig.AddFingerprint(serverConfig.Fingerprint{
		Name:         name,
		SerialNumber: serial,
		Fingerprint:  fingerprint, // This is the serial number for now
		AllowConnect: allow,
	})

	ac.clientLock.Lock()
	defer ac.clientLock.Unlock()
	ac.AcceptedClients = append(ac.AcceptedClients, ClientDetails{
		Name:        name,
		Serial:      serial,
		Fingerprint: fingerprint,
	})
}

func isCertificatePinned(cert *x509.Certificate) bool {
	return pinnedCertificates[cert.SerialNumber.String()]
}
