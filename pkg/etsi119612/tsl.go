package etsi119612

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"io"
	"net/http"
)

type Lang string

type TSL struct {
	StatusList TrustServiceStatusList `xml:"TrustServiceStatusList"`
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
	t := TSL{Source: url}

	err = xml.Unmarshal(bodyBytes, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (tsl *TSL) withTrustServices(cb func(int, *TSPType)) {
	for i, tsp := range tsl.StatusList.TslTrustServiceProviderList.TslTrustServiceProvider {
		cb(i, tsp)
	}
}

func (tsl *TSL) ToCertPool(policy *TSPServicePolicy) *x509.CertPool {
	pool := x509.NewCertPool()
	tsl.withTrustServices(func(i int, tsp *TSPType) {
		tsp.withCertificates(func(certData string) {
			data, err := base64.StdEncoding.DecodeString(certData)
			if err == nil {
				cert, err := x509.ParseCertificate(data)
				if err == nil {
					pool.AddCertWithConstraint(cert, func(chain []*x509.Certificate) error {
						return tsp.Validate(chain, policy)
					})
				} else {
					// TODO error logging
				}
			} else {
				//TODO error logging
			}
		})
	})
	return pool
}
