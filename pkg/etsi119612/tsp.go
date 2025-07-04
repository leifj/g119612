package etsi119612

import (
	"crypto/x509"
	"encoding/base64"
	"slices"

	log "github.com/sirupsen/logrus"
)

const ServiceStatusGranted string = "https://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted/"

// A struct representing configuration of the validation process. By default the ServiceStatus field
// contains a single element (ServiceStatusGranted) that represents the standardized value for indicating
// that the trust service provider is valid and granted access in the trust status list (ie not withdrawn).
// The ServiceTypeIdentifier is a list of allowed service types. When creating the CertPool for use in
// certificate validation the ServiceTypeIdentifier can be populated with a list of allowed types. If left
// empty this means every service type is allowed.
type TSPServicePolicy struct {
	ServiceTypeIdentifier []string
	ServiceStatus         []string
}

// A constant TSPServicePolicy instance that represents a standard policy with an empty ServiceTypeIdentifier array.
var (
	PolicyAll = NewTSPServicePolicy()
)

// Add an element to the ServiceTypeIdentifier array.
func (tc *TSPServicePolicy) AddServiceTypeIdentifier(sti string) {
	tc.ServiceTypeIdentifier = append(tc.ServiceTypeIdentifier, sti)
}

// Add an element to the ServiceStatus array. Note that adding to this array without first removing the standard "granted"
// element may not yield the expected results.
func (tc *TSPServicePolicy) AddServiceStatus(status string) {
	tc.ServiceStatus = append(tc.ServiceStatus, status)
}

// Create a standard TSPServicePolicy instance. Calling this creates the same object as the "PolicyAll" constant.
func NewTSPServicePolicy() *TSPServicePolicy {
	tc := TSPServicePolicy{ServiceTypeIdentifier: make([]string, 0), ServiceStatus: make([]string, 0)}
	tc.AddServiceStatus(ServiceStatusGranted)
	return &tc
}

// Cahe provided callback for all t all the X509 certificate data for the given Trust Service object.
func (svc *TSPServiceType) withCertificates(cb func(*x509.Certificate)) {
	if svc.TslServiceInformation.TslServiceDigitalIdentity != nil {
		for _, id := range svc.TslServiceInformation.TslServiceDigitalIdentity.DigitalId {
			if len(id.X509Certificate) > 0 {
				data, err := base64.StdEncoding.DecodeString(string(id.X509Certificate))
				if err == nil {
					cert, err := x509.ParseCertificate(data)
					if err == nil {
						cb(cert)
					} else {
						log.Errorf("g119612: [TSP: %s] Error parsing certificate: %s", FindByLanguage(svc.TslServiceInformation.ServiceName, "en", "Unknown"), err)
					}
				} else {
					log.Errorf("g119612: [TSP: %s] Error decoding certificate: %s", FindByLanguage(svc.TslServiceInformation.ServiceName, "en", "Unknown"), err)
				}
			}
		}
	}
}

// Checks a Trust Service for validity during certificate validation.
func (tsp *TSPType) Validate(svc *TSPServiceType, chain []*x509.Certificate, policy *TSPServicePolicy) error {

	if !slices.Contains(policy.ServiceStatus, svc.TslServiceInformation.TslServiceStatus) {
		return ErrInvalidStatus
	}

	if len(policy.ServiceTypeIdentifier) > 0 && !slices.Contains(policy.ServiceTypeIdentifier, svc.TslServiceInformation.TslServiceTypeIdentifier) {
		return ErrInvalidConstraints
	}

	return nil
}
