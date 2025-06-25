package etsi119612_test

import (
	"slices"
	"testing"

	"github.com/SUNET/g119612/pkg/etsi119612"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	defer gock.Off()
	gock.New("https://ewc-consortium.github.io").
		Get("/EWC-TL").
		Reply(200).
		File("./testdata/EWC-TL.xml")

	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	assert.NotNil(t, tsl)
	assert.Nil(t, err)
	assert.NotNil(t, tsl.StatusList)
	si := tsl.StatusList.TslSchemeInformation
	assert.NotNil(t, si)
	assert.Equal(t, si.TSLSequenceNumber, 1)
	assert.Equal(t, si.TSLSequenceNumber, 1)
	assert.Equal(t, *si.TslSchemeOperatorName.Name[0].XmlLangAttr, etsi119612.Lang("en"))
	assert.Equal(t, etsi119612.FindByLanguage(si.TslSchemeOperatorName, "en", "unknown"), "EWC Consortium")
	assert.Equal(t, etsi119612.FindByLanguage(si.TslSchemeOperatorName, "fr", "unknown 4711"), "unknown 4711")
}

func TestFetchMissingSchemeInfo(t *testing.T) {
	defer gock.Off()
	gock.New("https://ewc-consortium.github.io").
		Get("/EWC-TL").
		Reply(200).
		File("./testdata/EWC-TL-no-scheme-information.xml")

	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	assert.NotNil(t, tsl)
	assert.Nil(t, err)
	si := tsl.StatusList.TslSchemeInformation
	assert.Nil(t, si)
}

func TestFetchBrokenXML(t *testing.T) {
	defer gock.Off()
	gock.New("https://ewc-consortium.github.io").
		Get("/EWC-TL").
		Reply(200).
		File("./testdata/not-xml.xml")

	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	assert.Nil(t, tsl)
	assert.NotNil(t, err)
}

func TestFetchMissing(t *testing.T) {
	defer gock.Off()
	gock.New("https://example.com").
		Get("/missing").
		Reply(404)

	tsl, err := etsi119612.FetchTSL("https://example.com/missing")
	assert.Nil(t, tsl)
	assert.NotNil(t, err)
}

func TestFetchError(t *testing.T) {
	defer gock.Off()
	gock.New("https://example.com").
		Get("/bad").
		Reply(500)

	tsl, err := etsi119612.FetchTSL("https://example.com/bad")
	assert.Nil(t, tsl)
	assert.NotNil(t, err)
}

func TestFetchNotURL(t *testing.T) {
	tsl, err := etsi119612.FetchTSL("urn:not-an url")
	assert.Nil(t, tsl)
	assert.NotNil(t, err)
}

func TestCertPoolBadBase64(t *testing.T) {
	defer gock.Off()
	gock.New("https://ewc-consortium.github.io").
		Get("/EWC-TL").
		Reply(200).
		File("testdata/EWC-TL-bad-base64.xml")

	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	assert.NotNil(t, tsl)
	assert.Nil(t, err)
	pool := tsl.ToCertPool(etsi119612.PolicyAll)
	assert.NotNil(t, pool)
}

func TestCertPoolBadCert(t *testing.T) {
	defer gock.Off()
	gock.New("https://ewc-consortium.github.io").
		Get("/EWC-TL").
		Reply(200).
		File("testdata/EWC-TL-bad-cert.xml")

	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	assert.NotNil(t, tsl)
	assert.Nil(t, err)
	pool := tsl.ToCertPool(etsi119612.PolicyAll)
	assert.NotNil(t, pool)
}

func TestCertPool(t *testing.T) {
	defer gock.Off()
	gock.New("https://ewc-consortium.github.io").
		Get("/EWC-TL").
		Reply(200).
		File("testdata/EWC-TL.xml")

	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	assert.NotNil(t, tsl)
	assert.Nil(t, err)
	pool := tsl.ToCertPool(etsi119612.PolicyAll)
	assert.NotNil(t, pool)
}

func TestPolicy(t *testing.T) {
	p := etsi119612.NewTSPServicePolicy()
	assert.True(t, slices.ContainsFunc(p.ServiceStatus, func(s string) bool { return s == "https://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted/" }))
	assert.Equal(t, len(p.ServiceStatus), 1)
	p.AddServiceTypeIdentifier("urn:foo")
	assert.True(t, slices.ContainsFunc(p.ServiceTypeIdentifier, func(s string) bool { return s == "urn:foo" }))
	p.AddServiceStatus("urn:bar")
	assert.True(t, slices.ContainsFunc(p.ServiceStatus, func(s string) bool { return s == "urn:bar" }))
	assert.Equal(t, len(p.ServiceStatus), 2)
}
