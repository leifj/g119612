package etsi119612_test

import (
	"testing"

	"github.com/leifj/g119612/pkg/etsi119612"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	assert.NotNil(t, tsl)
	assert.Nil(t, err)
	assert.NotNil(t, tsl.StatusList)
	si := tsl.StatusList.TslSchemeInformation
	assert.NotNil(t, si)
	assert.Equal(t, si.TSLSequenceNumber, 1)
	assert.Equal(t, si.TSLSequenceNumber, 1)
}

func TestCertPool(t *testing.T) {
	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	assert.NotNil(t, tsl)
	assert.Nil(t, err)
	pool := tsl.ToCertPool(etsi119612.PolicyAll)
	assert.NotNil(t, pool)
}
