package etsi119612

import (
        "github.com/subchen/go-xmldom"
	"crypto/x509"
)

type TSPServicePolicy struct {
	ServiceTypeIdentifier []string
	ServiceStatus         []string
}

type TSPService struct {
        Node *xmldom.Node
}

func (tc *TSPServicePolicy) addServiceTypeIdentifier(sti string) {
        tc.ServiceTypeIdentifier = append(tc.ServiceTypeIdentifier,sti)
}

func (tc *TSPServicePolicy) addServiceStatus(status string) {
        tc.ServiceStatus = append(tc.ServiceStatus,status)
}


func NewTSPServicePolicy() *TSPServicePolicy {
	tc := TSPServicePolicy{ServiceTypeIdentifier: make([]string,0),ServiceStatus: make([]string,0)}
	tc.addServiceStatus("https://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted/")
	return &tc
}

var (
	PolicyAll = NewTSPServicePolicy()
)

func NewTSPService(node *xmldom.Node) *TSPService {
	t := TSPService{Node: node}
        return &t
}

func (*TSPService) Validate([]*x509.Certificate, *TSPServicePolicy) error {
	return nil //TBD
}
