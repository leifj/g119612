package etsi119612

import (
	"crypto/x509"
	"encoding/base64"
	"log"
)

type TSPServicePolicy struct {
	ServiceTypeIdentifier []string
	ServiceStatus         []string
}

func (tc *TSPServicePolicy) AddServiceTypeIdentifier(sti string) {
	tc.ServiceTypeIdentifier = append(tc.ServiceTypeIdentifier, sti)
}

func (tc *TSPServicePolicy) AddServiceStatus(status string) {
	tc.ServiceStatus = append(tc.ServiceStatus, status)
}

func NewTSPServicePolicy() *TSPServicePolicy {
	tc := TSPServicePolicy{ServiceTypeIdentifier: make([]string, 0), ServiceStatus: make([]string, 0)}
	tc.AddServiceStatus("https://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted/")
	return &tc
}

var (
	PolicyAll = NewTSPServicePolicy()
)

func (svc *TSPServiceType) withCertificates(cb func(*x509.Certificate)) {
	if svc.TslServiceInformation.TslServiceDigitalIdentity != nil {
		for _, id := range svc.TslServiceInformation.TslServiceDigitalIdentity.DigitalId {
			if len(id.X509Certificate) > 0 {
				data, err := base64.StdEncoding.DecodeString(string(id.X509Certificate))
				if err == nil {
					cert, err := x509.ParseCertificate(data)
					if err == nil {
						log.Printf("[TSP %s] Parsing certificate\n", FindByLanguage(svc.TslServiceInformation.ServiceName, "en", "Unknown"))
						cb(cert)
					} else {
						log.Printf("[TSP: %s] Error parsing certificate: %s", FindByLanguage(svc.TslServiceInformation.ServiceName, "en", "Unknown"), err)
					}
				} else {
					log.Printf("[TSP: %s] Error decoding certificate: %s", FindByLanguage(svc.TslServiceInformation.ServiceName, "en", "Unknown"), err)
				}
			}
		}
	}
}

func (*TSPServiceType) Validate([]*x509.Certificate, *TSPServicePolicy) error {
	return nil //TBD
}
