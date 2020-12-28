package shell

import (
	"reflect"
	"testing"

	"github.com/aliyun/saml2alibabacloud/pkg/alibabacloudconfig"
	"github.com/aliyun/saml2alibabacloud/pkg/cfg"
	"github.com/aliyun/saml2alibabacloud/pkg/flags"
)

func TestBuildEnvVars(t *testing.T) {
	account := &cfg.IDPAccount{
		Profile: "saml",
	}
	alibabacloudCreds := &alibabacloudconfig.AliCloudCredentials{
		AliCloudAccessKey:     "123",
		AliCloudSecretKey:     "345",
		AliCloudSecurityToken: "567",
		AliCloudSessionToken:  "567",
	}

	tests := []struct {
		name  string
		flags *flags.LoginExecFlags
		want  []string
	}{
		{
			name:  "build-env",
			flags: &flags.LoginExecFlags{},
			want: []string{
				"ALICLOUD_ASSUME_ROLE_SESSION_NAME=567",
				"ALICLOUD_SECURITY_TOKEN=567",
				"ALICLOUD_ACCESS_KEY=123",
				"ALICLOUD_SECRET_KEY=345",
			},
		},
		{
			name: "build-env-with-profile",
			flags: &flags.LoginExecFlags{
				ExecProfile: "testing",
			},
			want: []string{
				"ALICLOUD_ASSUME_ROLE_SESSION_NAME=567",
				"ALICLOUD_SECURITY_TOKEN=567",
				"ALICLOUD_ACCESS_KEY=123",
				"ALICLOUD_SECRET_KEY=345",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildEnvVars(alibabacloudCreds, account, tt.flags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildEnvVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
