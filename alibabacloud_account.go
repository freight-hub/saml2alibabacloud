package saml2alibabacloud

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// AlibabaCloudAccount holds the AlibabaCloud account name and roles
type AlibabaCloudAccount struct {
	Name  string
	Roles []*RamRole
}

// ParseAlibabaCloudAccounts extract the AlibabaCloud accounts from the saml assertion
func ParseAlibabaCloudAccounts(audience string, samlAssertion string) ([]*AlibabaCloudAccount, error) {
	res, err := http.PostForm(audience, url.Values{"SAMLResponse": {samlAssertion}})
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving AlibabaCloud login form")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving AlibabaCloud login body")
	}

	return ExtractAlibabaCloudAccounts(data)
}

// ExtractAlibabaCloudAccounts extract the accounts from the AlibabaCloud login html page
func ExtractAlibabaCloudAccounts(data []byte) ([]*AlibabaCloudAccount, error) {
	accounts := []*AlibabaCloudAccount{}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to build document from response")
	}

	doc.Find("#samlRoleForm > div.form-group > div.col-sm-4 > label").Each(func(i int, s *goquery.Selection) {
		account := new(AlibabaCloudAccount)
		name := s.Text()
		parts := strings.Split(name, ":")
		account.Name = strings.TrimSpace(parts[1])
		s.Parent().Parent().Next().Find("input[name='roleAttribute']").Each(func(i int, s *goquery.Selection) {
			role := new(RamRole)
			label := s.Parent()
			role.Name = strings.TrimSpace(label.Text())
			value, _ := s.Attr("value")
			parts = strings.Split(value, ",")
			role.RoleARN = strings.TrimSpace(parts[0])
			role.PrincipalARN = strings.TrimSpace(parts[1])
			account.Roles = append(account.Roles, role)
		})
		accounts = append(accounts, account)
	})

	return accounts, nil
}

// AssignPrincipals assign principal from roles
func AssignPrincipals(ramRoles []*RamRole, alibabacloudAccounts []*AlibabaCloudAccount) {

	principalARNs := make(map[string]string)
	for _, ramRole := range ramRoles {
		principalARNs[ramRole.RoleARN] = ramRole.PrincipalARN
	}

	for _, account := range alibabacloudAccounts {
		for _, ramRole := range account.Roles {
			ramRole.PrincipalARN = principalARNs[ramRole.RoleARN]
		}
	}

}

// LocateRole locate role by name
func LocateRole(ramRoles []*RamRole, roleName string) (*RamRole, error) {
	for _, ramRole := range ramRoles {
		if ramRole.RoleARN == roleName {
			return ramRole, nil
		}
	}

	return nil, fmt.Errorf("Supplied RoleArn not found in saml assertion: %s", roleName)
}
