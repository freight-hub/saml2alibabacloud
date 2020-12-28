package saml2alibabacloud

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractAlibabaCloudAccounts(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/saml.html")
	assert.Nil(t, err)

	accounts, err := ExtractAlibabaCloudAccounts(data)
	assert.Nil(t, err)
	assert.Len(t, accounts, 2)

	account := accounts[0]
	assert.Equal(t, account.Name, "Account: account-alias (000000000001)")

	assert.Len(t, account.Roles, 2)
	role := account.Roles[0]
	assert.Equal(t, role.RoleARN, "acs:ram::000000000001:role/Development")
	assert.Equal(t, role.Name, "Development")
	role = account.Roles[1]
	assert.Equal(t, role.RoleARN, "acs:ram::000000000001:role/Production")
	assert.Equal(t, role.Name, "Production")

	account = accounts[1]
	assert.Equal(t, account.Name, "Account: 000000000002")

	assert.Len(t, account.Roles, 1)
	role = account.Roles[0]
	assert.Equal(t, role.RoleARN, "acs:ram::000000000002:role/Production")
	assert.Equal(t, role.Name, "Production")
}

func TestAssignPrincipals(t *testing.T) {
	ramRoles := []*RamRole{
		{
			PrincipalARN: "acs:ram::000000000001:saml-provider/test-idp",
			RoleARN:      "acs:ram::000000000001:role/Development",
		},
	}

	alibabacloudAccounts := []*AlibabaCloudAccount{
		{
			Roles: []*RamRole{
				{
					RoleARN: "acs:ram::000000000001:role/Development",
				},
			},
		},
	}

	AssignPrincipals(ramRoles, alibabacloudAccounts)

	assert.Equal(t, "acs:ram::000000000001:saml-provider/test-idp", alibabacloudAccounts[0].Roles[0].PrincipalARN)
}

func TestLocateRole(t *testing.T) {
	ramRoles := []*RamRole{
		{
			PrincipalARN: "acs:ram::000000000001:saml-provider/test-idp",
			RoleARN:      "acs:ram::000000000001:role/Development",
		},
		{
			PrincipalARN: "acs:ram::000000000002:saml-provider/test-idp",
			RoleARN:      "acs:ram::000000000002:role/Development",
		},
	}

	role, err := LocateRole(ramRoles, "acs:ram::000000000001:role/Development")

	assert.Empty(t, err)

	assert.Equal(t, "acs:ram::000000000001:role/Development", role.RoleARN)
}
