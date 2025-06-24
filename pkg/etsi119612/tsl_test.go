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
	si := tsl.StatusList.TslSchemeInformation
	assert.NotNil(t, si)
	assert.Equal(t, si.TSLVersionIdentifier, 1)
	assert.Equal(t, si.TSLSequenceNumber, 1)
}
