package saml2alibabacloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRoles(t *testing.T) {

	roles := []string{
		"acs:ram::456456456456:saml-provider/example-idp,acs:ram::456456456456:role/admin",
		"acs:ram::456456456456:role/admin,acs:ram::456456456456:saml-provider/example-idp",
	}

	ramRoles, err := ParseRamRoles(roles)

	assert.Nil(t, err)
	assert.Len(t, ramRoles, 2)

	for _, ramRole := range ramRoles {
		assert.Equal(t, "acs:ram::456456456456:saml-provider/example-idp", ramRole.PrincipalARN)
		assert.Equal(t, "acs:ram::456456456456:role/admin", ramRole.RoleARN)
	}

	roles = []string{""}
	ramRoles, err = ParseRamRoles(roles)

	assert.NotNil(t, err)
	assert.Nil(t, ramRoles)

}
