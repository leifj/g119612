package etsi119612_test

import (
	"testing"
	"os"
	"github.com/SUNET/g119612/pkg/etsi119612"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"fmt"
	"encoding/base64"
	"crypto/x509"
)

type JWTCertBundle struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	X5c []string `json:"x5c"`

}

func TestLeafRootCertVerificationSuccess(t *testing.T) {
	header_mock, err:= os.ReadFile("./testdata/x5c-test-root-leaf.json")
	if err != nil {
		t.Fatalf("Failed while reading json: %v", err)
	}
	assert.NoError(t, err)

    assert.NotEmpty(t, header_mock) 
	var jwt JWTCertBundle
	err =json.Unmarshal(header_mock, &jwt)
	if err !=nil{
		t.Fatalf("Failed updating jwt bundle")
	}
	assert.NotEmpty(t, jwt)
	assert.NotEmpty(t, jwt.Alg) 
	assert.NotEmpty(t, jwt.Typ)
	assert.NotEmpty(t, jwt.X5c)
	defer gock.Off()
	gock.New("https://ewc-consortium.github.io").
		Get("/EWC-TL").
		Reply(200).
		File("testdata/test-trust-list.xml")
	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	
	fmt.Println(tsl)

	
	pool := tsl.ToCertPool(etsi119612.PolicyAll)
	fmt.Println("Number of trusted roots:", len(pool.Subjects()))
	fmt.Println(pool)
	assert.NotNil(t, pool)
	leafDER, err:= base64.StdEncoding.DecodeString(jwt.X5c[0])
	assert.NoError(t,err)
	leafCert, err := x509.ParseCertificate(leafDER)
	_, err = leafCert.Verify(x509.VerifyOptions{
    Roots: pool})
	if err !=nil {
		t.Errorf("Leaf certificate verification failed %v", err)
	} else {
		fmt.Println("Leaf certificate verified against root")
	}
}

func TestLeafIntermediateRootCertVerificationSuccess(t *testing.T) {
	header_mock, err:= os.ReadFile("./testdata/x5c-test.json")
	if err != nil {
		t.Fatalf("Failed while reading json: %v", err)
	}
	assert.NoError(t, err)

    assert.NotEmpty(t, header_mock) 
	var jwt JWTCertBundle
	err =json.Unmarshal(header_mock, &jwt)
	if err !=nil{
		t.Fatalf("Failed updating jwt bundle")
	}
	assert.NotEmpty(t, jwt)
	assert.NotEmpty(t, jwt.Alg) 
	assert.NotEmpty(t, jwt.Typ)
	assert.NotEmpty(t, jwt.X5c)
	defer gock.Off()
	gock.New("https://ewc-consortium.github.io").
		Get("/EWC-TL").
		Reply(200).
		File("testdata/test-trust-list.xml")
	tsl, err := etsi119612.FetchTSL("https://ewc-consortium.github.io/ewc-trust-list/EWC-TL")
	
	fmt.Println(tsl)

	
	pool := tsl.ToCertPool(etsi119612.PolicyAll)
	fmt.Println("Number of trusted roots:", len(pool.Subjects()))
	fmt.Println(pool)
	assert.NotNil(t, pool)
	leafDER, err:= base64.StdEncoding.DecodeString(jwt.X5c[0])
	assert.NoError(t,err)
	leafCert, err := x509.ParseCertificate(leafDER)
	//adding the intermediate certificate from x5c
	interDER, err := base64.StdEncoding.DecodeString(jwt.X5c[1])
    assert.NoError(t, err, "Failed to decode intermediate")

    interCert, err := x509.ParseCertificate(interDER)
	assert.NoError(t, err, "Failed to parse intermediate certificate")

	interPool := x509.NewCertPool()
	interPool.AddCert(interCert)
	assert.NotEmpty(t, interPool.Subjects(), "Intermediate cert pool is empty")
	_, err = leafCert.Verify(x509.VerifyOptions{
    Roots: pool, Intermediates:interPool})
	if err !=nil {
		t.Errorf("Leaf certificate verification failed %v", err)
	} else {
		fmt.Println("Leaf certificate verified against root")
	}
}
