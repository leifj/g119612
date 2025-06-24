package etsi119612

import (
	"slices"
	"github.com/subchen/go-xmldom"
 	"net/http"
	"crypto/x509"
	"encoding/base64"
)

type TSL struct {
        Doc *xmldom.Document
        Source string
}

func FetchTSL(url string) (*TSL,error) {
        resp, err := http.Get(url)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()
        doc := xmldom.Must(xmldom.Parse(resp.Body))
        t := TSL{Source: url}
        t.Doc = doc
        return &t,nil
}

func (tsl *TSL) findTrustServiceProviderList() *xmldom.Node {
	return tsl.Doc.Root.FindOneByName("TrustServiceProviderList")
}

func (tsl *TSL) withTSP(cb func (int, *xmldom.Node)) {
	tsl.Doc.Root.QueryEach("//TrustServiceProvider",cb)
}

func (tsl *TSL) withTrustServices(cb func (int, *xmldom.Node)) {
	tsl.Doc.Root.QueryEach("//TSPService", cb)
}

func (tsl *TSL) FindServicesBySKI(skiVal string) []*xmldom.Node {
	nodes := make([]*xmldom.Node,0)
	tsl.withTSP(func (i int, n *xmldom.Node) {
		skis := n.FindByName("X509SKI")
		idx := slices.IndexFunc(skis,func(ski *xmldom.Node) bool { return ski.Text == skiVal })
		if (idx != -1) {
			nodes = append(nodes,n)
                }
	})
	return nodes
}

func (tsl *TSL) ToCertPool(policy *TSPServicePolicy) *x509.CertPool {
	pool := x509.NewCertPool()
	tsl.withTrustServices(func (i int, tspNode *xmldom.Node) {
		certNodes := tspNode.FindByName("X509Certificate")
		for _,n := range certNodes {
			tsp := NewTSPService(tspNode)
			data,err := base64.StdEncoding.DecodeString(n.Text)
			if (err == nil) {
				cert,err := x509.ParseCertificate(data)
				if (err == nil) {
					pool.AddCertWithConstraint(cert, func(chain []*x509.Certificate) error {
						return tsp.Validate(chain,policy)
					})
				} else {
					// TODO error logging
				}
			} else {
				//TODO error logging
			}
		}
	});
	return pool
}
