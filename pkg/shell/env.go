package shell

import (
	"fmt"

	"github.com/aliyun/saml2alibabacloud/pkg/alibabacloudconfig"
	"github.com/aliyun/saml2alibabacloud/pkg/cfg"
	"github.com/aliyun/saml2alibabacloud/pkg/flags"
)

// BuildEnvVars build an array of env vars in the format required for exec
func BuildEnvVars(alibabacloudCreds *alibabacloudconfig.AliCloudCredentials, account *cfg.IDPAccount, execFlags *flags.LoginExecFlags) []string {

	environmentVars := []string{
		fmt.Sprintf("ALICLOUD_ASSUME_ROLE_SESSION_NAME=%s", alibabacloudCreds.AliCloudSessionToken),
		fmt.Sprintf("ALICLOUD_SECURITY_TOKEN=%s", alibabacloudCreds.AliCloudSecurityToken),
		fmt.Sprintf("ALICLOUD_ACCESS_KEY=%s", alibabacloudCreds.AliCloudAccessKey),
		fmt.Sprintf("ALICLOUD_SECRET_KEY=%s", alibabacloudCreds.AliCloudSecretKey),
	}

	if execFlags.ExecProfile == "" {
		// Only set profile env vars if we haven't already assumed a role via a profile
		environmentVars = append(environmentVars, fmt.Sprintf("ALICLOUD_PROFILE=%s", account.Profile))
	}
	return environmentVars
}
