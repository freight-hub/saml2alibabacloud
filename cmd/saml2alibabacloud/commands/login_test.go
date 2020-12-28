package commands

import (
	"testing"

	saml2alibabacloud "github.com/aliyun/saml2alibabacloud"
	"github.com/aliyun/saml2alibabacloud/pkg/cfg"
	"github.com/aliyun/saml2alibabacloud/pkg/creds"
	"github.com/aliyun/saml2alibabacloud/pkg/flags"
	"github.com/stretchr/testify/assert"
)

func TestResolveLoginDetailsWithFlags(t *testing.T) {

	commonFlags := &flags.CommonFlags{URL: "https://id.example.com", Username: "ziying", Password: "alibabacloud", MFAToken: "123456", SkipPrompt: true}
	loginFlags := &flags.LoginExecFlags{CommonFlags: commonFlags}

	idpa := &cfg.IDPAccount{
		URL:      "https://id.example.com",
		MFA:      "none",
		Provider: "Ping",
		Username: "ziying",
	}
	loginDetails, err := resolveLoginDetails(idpa, loginFlags)

	assert.Empty(t, err)
	assert.Equal(t, &creds.LoginDetails{Username: "ziying", Password: "alibabacloud", URL: "https://id.example.com", MFAToken: "123456"}, loginDetails)
}

func TestResolveRoleSingleEntry(t *testing.T) {

	adminRole := &saml2alibabacloud.RamRole{
		Name:         "admin",
		RoleARN:      "acs:ram::1234567890:role/ali-cloudadmin-master",
		PrincipalARN: "acs:ram::1234567890:role/ali-cloudadmin-master,acs:ram::1234567890:saml-provider/example-idp",
	}

	alibabacloudRoles := []*saml2alibabacloud.RamRole{
		adminRole,
	}

	got, err := resolveRole(alibabacloudRoles, "", cfg.NewIDPAccount())
	assert.Empty(t, err)
	assert.Equal(t, got, adminRole)
}
