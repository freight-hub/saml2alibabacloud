package commands

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"os"

	saml2alibabacloud "github.com/aliyun/saml2alibabacloud"
	"github.com/aliyun/saml2alibabacloud/helper/credentials"
	"github.com/aliyun/saml2alibabacloud/pkg/flags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ListRoles will list available role ARNs
func ListRoles(loginFlags *flags.LoginExecFlags) error {

	logger := logrus.WithField("command", "list")

	account, err := buildIdpAccount(loginFlags)
	if err != nil {
		return errors.Wrap(err, "error building login details")
	}

	loginDetails, err := resolveLoginDetails(account, loginFlags)
	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	err = loginDetails.Validate()
	if err != nil {
		return errors.Wrap(err, "error validating login details")
	}

	logger.WithField("idpAccount", account).Debug("building provider")

	provider, err := saml2alibabacloud.NewSAMLClient(account)
	if err != nil {
		return errors.Wrap(err, "error building IdP client")
	}

	samlAssertion, err := provider.Authenticate(loginDetails)
	if err != nil {
		return errors.Wrap(err, "error authenticating to IdP")
	}

	if samlAssertion == "" {
		log.Println("Response did not contain a valid SAML assertion")
		log.Println("Please check your username and password is correct")
		log.Println("To see the output follow the instructions in https://github.com/aliyun/saml2alibabacloud#debugging-issues-with-idps")
		os.Exit(1)
	}

	if !loginFlags.CommonFlags.DisableKeychain {
		err = credentials.SaveCredentials(loginDetails.URL, loginDetails.Username, loginDetails.Password)
		if err != nil {
			return errors.Wrap(err, "error storing password in keychain")
		}
	}

	data, err := b64.StdEncoding.DecodeString(samlAssertion)
	if err != nil {
		return errors.Wrap(err, "error decoding saml assertion")
	}

	roles, err := saml2alibabacloud.ExtractRamRoles(data)
	if err != nil {
		return errors.Wrap(err, "error parsing AlibabaCloud roles")
	}

	if len(roles) == 0 {
		log.Println("No roles to assume")
		os.Exit(1)
	}

	alibabacloudRoles, err := saml2alibabacloud.ParseRamRoles(roles)
	if err != nil {
		return errors.Wrap(err, "error parsing AlibabaCloud roles")
	}

	if err := listRoles(alibabacloudRoles, samlAssertion, loginFlags); err != nil {
		return errors.Wrap(err, "Failed to list roles")
	}

	return nil
}

func listRoles(alibabacloudRoles []*saml2alibabacloud.RamRole, samlAssertion string, loginFlags *flags.LoginExecFlags) error {
	if len(alibabacloudRoles) == 1 {
		log.Println("")
		log.Println("Only one role to assume. Will be automatically assumed on login")
		log.Println(alibabacloudRoles[0].RoleARN)
		return nil
	} else if len(alibabacloudRoles) == 0 {
		return errors.New("no roles available")
	}

	samlAssertionData, err := b64.StdEncoding.DecodeString(samlAssertion)
	if err != nil {
		return errors.Wrap(err, "error decoding saml assertion")
	}

	aud, err := saml2alibabacloud.ExtractDestinationURL(samlAssertionData)
	if err != nil {
		return errors.Wrap(err, "error parsing destination url")
	}

	alibabacloudAccounts, err := saml2alibabacloud.ParseAlibabaCloudAccounts(aud, samlAssertion)
	if err != nil {
		return errors.Wrap(err, "error parsing AlibabaCloud role accounts")
	}

	saml2alibabacloud.AssignPrincipals(alibabacloudRoles, alibabacloudAccounts)

	log.Println("")
	for _, account := range alibabacloudAccounts {
		fmt.Println(account.Name)
		for _, role := range account.Roles {
			fmt.Println(role.RoleARN)
		}
		fmt.Println("")
	}

	return nil
}
