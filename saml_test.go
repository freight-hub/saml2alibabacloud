package saml2alibabacloud

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractRamRoles(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/assertion.xml")
	assert.Nil(t, err)

	roles, err := ExtractRamRoles(data)
	assert.Nil(t, err)
	assert.Len(t, roles, 2)
}

func TestExtractSessionDuration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/assertion.xml")
	assert.Nil(t, err)

	duration, err := ExtractSessionDuration(data)
	assert.Nil(t, err)
	assert.Equal(t, int64(28800), duration)
}

func TestExtractDestinationURL(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/assertion.xml")
	assert.Nil(t, err)

	destination, err := ExtractDestinationURL(data)
	assert.Nil(t, err)
	assert.Equal(t, "https://signin.aliyun.com/saml-role/sso", destination)
}

func TestExtractDestinationURL2(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/assertion_no_destination.xml")
	assert.Nil(t, err)

	destination, err := ExtractDestinationURL(data)
	assert.Nil(t, err)
	assert.Equal(t, "https://signin.aliyun.com/saml-role/sso", destination)
}
