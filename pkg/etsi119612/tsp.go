package etsi119612

import (
	"crypto/x509"
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

func (tsp *TSPType) withCertificates(cb func(string)) {
	for _, svc := range tsp.TslTSPServices.TslTSPService {
		for _, id := range svc.TslServiceInformation.TslServiceDigitalIdentity.DigitalId {
			if len(id.X509Certificate) > 0 {
				cb(id.X509Certificate)
			}
		}
	}
}

func (*TSPType) Validate([]*x509.Certificate, *TSPServicePolicy) error {
	return nil //TBD
}
