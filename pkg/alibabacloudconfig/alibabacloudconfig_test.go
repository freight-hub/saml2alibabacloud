package alibabacloudconfig

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestUpdateSamlConfig(t *testing.T) {
	os.Remove(".credentials")

	logrus.SetLevel(logrus.DebugLevel)

	sharedCreds := &CredentialsProvider{".credentials", "saml"}

	exist, err := sharedCreds.CredsExists()
	assert.Nil(t, err)
	assert.True(t, exist)

	alibabacloudCreds := &AliCloudCredentials{
		AliCloudAccessKey:     "testid",
		AliCloudSecretKey:     "testsecret",
		AliCloudSessionToken:  "testtoken",
		AliCloudSecurityToken: "testtoken",
	}

	err = sharedCreds.Save(alibabacloudCreds)
	assert.Nil(t, err)

	profile, err := sharedCreds.Load()
	assert.Nil(t, err)
	assert.Equal(t, "testid", profile.AliCloudAccessKey)
	assert.Equal(t, "testsecret", profile.AliCloudSecretKey)
	assert.Equal(t, "testtoken", profile.AliCloudSecurityToken)

	os.Remove(".credentials")
}
