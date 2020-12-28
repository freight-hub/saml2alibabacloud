package commands

import (
	"log"
	"os"
	"text/template"

	"github.com/aliyun/saml2alibabacloud/pkg/alibabacloudconfig"
	"github.com/aliyun/saml2alibabacloud/pkg/flags"
	"github.com/pkg/errors"
)

const bashTmpl = `export ALIBABA_CLOUD_ACCESS_KEY_ID="{{ .AliCloudAccessKey }}"
export ALIBABA_CLOUD_ACCESS_KEY_SECRET="{{ .AliCloudSecretKey }}"
export ALIBABA_CLOUD_SESSION_TOKEN="{{ .AliCloudSessionToken }}"
export ALIBABA_CLOUD_SECURITY_TOKEN="{{ .AliCloudSecurityToken }}"
export ALICLOUD_ACCESS_KEY="{{ .AliCloudAccessKey }}"
export ALICLOUD_SECRET_KEY="{{ .AliCloudSecretKey }}"
export ALICLOUD_SECURITY_TOKEN="{{ .AliCloudSecurityToken }}"
export ALICLOUD_PROFILE="{{ .ProfileName }}"
export SAML2ALIBABA_CLOUD_PROFILE="{{ .ProfileName }}"
`

const fishTmpl = `set -gx ALIBABA_CLOUD_ACCESS_KEY_ID {{ .AliCloudAccessKey }}
set -gx ALIBABA_CLOUD_ACCESS_KEY_SECRET {{ .AliCloudSecretKey }}
set -gx ALIBABA_CLOUD_SESSION_TOKEN {{ .AliCloudSessionToken }}
set -gx ALIBABA_CLOUD_SECURITY_TOKEN {{ .AliCloudSecurityToken }}
set -gx ALICLOUD_ACCESS_KEY {{ .AliCloudAccessKey }}
set -gx ALICLOUD_SECRET_KEY {{ .AliCloudSecretKey }}
set -gx ALICLOUD_SECURITY_TOKEN {{ .AliCloudSecurityToken }}
set -gx ALICLOUD_PROFILE {{ .ProfileName }}
set -gx SAML2ALIBABA_CLOUD_PROFILE {{ .ProfileName }}
`

const powershellTmpl = `$env:ALIBABA_CLOUD_ACCESS_KEY_ID='{{ .AliCloudAccessKey }}'
$env:ALIBABA_CLOUD_ACCESS_KEY_SECRET='{{ .AliCloudSecretKey }}'
$env:ALIBABA_CLOUD_SESSION_TOKEN='{{ .AliCloudSessionToken }}'
$env:ALIBABA_CLOUD_SECURITY_TOKEN='{{ .AliCloudSecurityToken }}'
$env:ALICLOUD_ACCESS_KEY='{{ .AliCloudAccessKey }}'
$env:ALICLOUD_SECRET_KEY='{{ .AliCloudSecretKey }}'
$env:ALICLOUD_SECURITY_TOKEN='{{ .AliCloudSecurityToken }}'
$env:ALICLOUD_PROFILE='{{ .ProfileName }}'
$env:SAML2ALIBABA_CLOUD_PROFILE='{{ .ProfileName }}'
`

// Script will emit a bash script that will export environment variables
func Script(execFlags *flags.LoginExecFlags, shell string) error {
	account, err := buildIdpAccount(execFlags)
	if err != nil {
		return errors.Wrap(err, "error building login details")
	}

	sharedCreds := alibabacloudconfig.NewSharedCredentials(account.Profile)

	// this checks if the credentials file has been created yet
	// can only really be triggered if saml2alibabacloud exec is run on a new
	// system prior to creating $HOME/.aliyun
	exist, err := sharedCreds.CredsExists()
	if err != nil {
		return errors.Wrap(err, "error loading credentials")
	}
	if !exist {
		log.Println("unable to load credentials, login required to create them")
		return nil
	}

	alibabacloudCreds, err := sharedCreds.Load()
	if err != nil {
		return errors.Wrap(err, "error loading credentials")
	}

	// annoymous struct to pass to template
	data := struct {
		ProfileName string
		*alibabacloudconfig.AliCloudCredentials
	}{
		account.Profile,
		alibabacloudCreds,
	}

	err = buildTmpl(shell, data)
	if err != nil {
		return errors.Wrap(err, "error generating template")
	}

	return nil
}

func buildTmpl(shell string, data interface{}) error {
	t := template.New("envvar_script")

	var err error

	switch shell {
	case "bash":
		t, err = t.Parse(bashTmpl)
	case "powershell":
		t, err = t.Parse(powershellTmpl)
	case "fish":
		t, err = t.Parse(fishTmpl)
	}

	if err != nil {
		return err
	}
	// this is still written to stdout as per convention
	return t.Execute(os.Stdout, data)
}
