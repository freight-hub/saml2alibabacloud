package cfg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const throwAwayConfig = "example/saml2alibabacloud.test.ini"

func TestNewConfigManagerNew(t *testing.T) {

	cfgm, err := NewConfigManager("example/saml2alibabacloud.ini")
	require.Nil(t, err)

	require.NotNil(t, cfgm)
}

func TestNewConfigManagerLoad(t *testing.T) {

	cfgm, err := NewConfigManager("example/saml2alibabacloud.ini")
	require.Nil(t, err)

	require.NotNil(t, cfgm)

	idpAccount, err := cfgm.LoadIDPAccount("test123")
	require.Nil(t, err)
	require.Equal(t, &IDPAccount{
		URL:             "https://id.whatever.com",
		Username:        "abc@whatever.com",
		Provider:        "keycloak",
		MFA:             "sms",
		AlibabaCloudURN: DefaultAlibabaCloudURN,
		SessionDuration: 3600,
		Profile:         "saml",
	}, idpAccount)

	idpAccount, err = cfgm.LoadIDPAccount("")
	require.Nil(t, err)
	require.Equal(t, &IDPAccount{
		AlibabaCloudURN: DefaultAlibabaCloudURN,
		SessionDuration: 3600,
		Profile:         "saml",
	}, idpAccount)
}

func TestNewConfigManagerSave(t *testing.T) {

	cfgm, err := NewConfigManager(throwAwayConfig)
	require.Nil(t, err)

	err = cfgm.SaveIDPAccount("testing2", &IDPAccount{
		URL:      "https://id.whatever.com",
		MFA:      "none",
		Provider: "keycloak",
		Username: "abc@whatever.com",
		Profile:  "saml",
	})
	require.Nil(t, err)
	idpAccount, err := cfgm.LoadIDPAccount("testing2")
	require.Nil(t, err)
	require.Equal(t, &IDPAccount{
		URL:             "https://id.whatever.com",
		Username:        "abc@whatever.com",
		Provider:        "keycloak",
		MFA:             "none",
		AlibabaCloudURN: DefaultAlibabaCloudURN,
		Profile:         "saml",
	}, idpAccount)

	os.Remove(throwAwayConfig)

}
