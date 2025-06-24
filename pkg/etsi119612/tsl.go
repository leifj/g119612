package etsi119612

import (
	"bytes"
	"crypto/x509"
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

type TSL struct {
	StatusList TrustStatusListType `xml:"tsl:TrustServiceStatusList"`
	Source     string
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func FetchTSL(url string) (*TSL, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes := StreamToByte(resp.Body)
	t := TSL{Source: url, StatusList: TrustStatusListType{}}
	log.Printf("Fetched %d bytes\n", len(bodyBytes))

	err = xml.Unmarshal(bodyBytes, &t.StatusList)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (tsl *TSL) withTrustServices(cb func(*TSPServiceType)) {
	for _, tsp := range tsl.StatusList.TslTrustServiceProviderList.TslTrustServiceProvider {
		if tsp != nil {
			log.Printf("TSP: %+v\n", *tsp.TslTSPInformation.TSPName.Name[0].NonEmptyNormalizedString)
			for _, svc := range tsp.TslTSPServices.TslTSPService {
				log.Printf("SVC: %+v\n", *svc.TslServiceInformation.ServiceName.Name[0].NonEmptyNormalizedString)
				cb(svc)
			}
		}
	}
}

func (tsl *TSL) ToCertPool(policy *TSPServicePolicy) *x509.CertPool {
	pool := x509.NewCertPool()
	tsl.withTrustServices(func(svc *TSPServiceType) {
		svc.withCertificates(func(cert *x509.Certificate) {
			pool.AddCertWithConstraint(cert, func(chain []*x509.Certificate) error {
				return svc.Validate(chain, policy)
			})
		})
	})
	return pool
}
