# saml2alicloud [![GitHub Actions status](https://github.com/aliyun/saml2alibabacloud/workflows/Go/badge.svg?branch=master)](https://github.com/aliyun/saml2alibabacloud/actions?query=workflow%3AGo) [![Build status - Windows](https://ci.appveyor.com/api/projects/status/ptpi18kci16o4i82/branch/master?svg=true)](https://ci.appveyor.com/project/davidobrien1985/saml2aws/branch/master)

CLI tool which enables you to login and retrieve [AlibabaCloud](https://www.aliyun.com/) temporary credentials using 
with [ADFS](https://msdn.microsoft.com/en-us/library/bb897402.aspx) or [PingFederate](https://www.pingidentity.com/en/products/pingfederate.html) Identity Providers.

This is based on [
saml2aws](https://github.com/Versent/saml2aws). Great thanks to Versent's work.

The process goes something like this:

* Setup an account alias, either using the default or given a name
* Prompt user for credentials
* Log in to Identity Provider using form based authentication
* Build a SAML assertion containing AlibabaCloud RAM roles
* Exchange the role and SAML assertion with [AlibabaCloud STS service](https://www.alibabacloud.com/help/doc-detail/109979.htm) to get a temporary set of credentials
* Save these credentials to an AlibabaCloud CLI profile named "saml"

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Requirements](#requirements)
- [Caveats](#caveats)
- [Install](#install)
    - [OSX](#osx)
    - [Windows](#windows)
    - [Linux](#linux)
- [Dependency Setup](#dependency-setup)
- [Usage](#usage)
    - [`saml2alicloud script`](#saml2alibabacloud-script)
    - [Configuring IDP Accounts](#configuring-idp-accounts)
- [Example](#example)
- [Advanced Configuration](#advanced-configuration)
    - [Dev Account Setup](#dev-account-setup)
    - [Test Account Setup](#test-account-setup)
- [Building](#building)
- [Environment vars](#environment-vars)
- [Provider Specific Documentation](#provider-specific-documentation)

## Requirements

* One of the supported Identity Providers
  * ADFS (2.x or 3.x)
  * [AzureAD](doc/provider/aad/README.md)
  * PingFederate + PingId
  * [Okta](pkg/provider/okta/README.md)
  * KeyCloak + (TOTP)
  * [Google Apps](pkg/provider/googleapps/README.md)
  * [Shibboleth](pkg/provider/shibboleth/README.md)
  * [F5APM](pkg/provider/f5apm/README.md)
  * [Akamai](pkg/provider/akamai/README.md)
  * OneLogin
  * NetIQ
* AlibabaCloud SAML Provider configured

## Caveats

Aside from Okta, most of the providers in this project are using screen scraping to log users into SAML, this isn't ideal and hopefully vendors make this easier in the future. In addition to this there are some things you need to know:

1. AlibabaCloud defaults to session tokens being issued with a duration of up to 3600 seconds (1 hour), this can now be configured as [Set the maximum session duration for a RAM role](https://www.alibabacloud.com/help/doc-detail/166256.htm).
2. Every SAML provider is different, the login process, MFA support is pluggable and therefore some work may be needed to integrate with your identity server

## Install

```
$ CURRENT_VERSION=0.0.1
$ wget https://github.com/aliyun/saml2alibabacloud/releases/download/v${CURRENT_VERSION}/saml2alibabacloud_${CURRENT_VERSION}_linux_amd64.tar.gz
$ tar -xzvf saml2alibabacloud_${CURRENT_VERSION}_linux_amd64.tar.gz -C /usr/local/bin
$ chmod u+x /usr/local/bin/saml2alibabacloud
$ saml2alibabacloud --version
```
**Note**: You will need to logout of your current user session or force a bash reload for `saml2alibabacloud` to be useable after following the above steps.

e.g. `exec -l bash`

#### [Void Linux](https://voidlinux.org/)

If you are on Void Linux you can use xbps to install the saml2alibabacloud package!

```
xbps-install saml2alibabacloud
```

## Dependency Setup

Install the AlibabaCloud CLI [see](alibabacloud.com/help/doc-detail/121544.htm)

## Usage

```
usage: saml2alibabacloud [<flags>] <command> [<args> ...]

A command line tool to help with SAML access to the AlibabaCloud STS service.

Flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
      --version                Show application version.
      --verbose                Enable verbose logging
  -i, --provider=PROVIDER      This flag is obsolete. See: https://github.com/aliyun/saml2alibabacloud#configuring-idp-accounts
  -a, --idp-account="default"  The name of the configured IDP account. (env: SAML2ALIBABACLOUD_IDP_ACCOUNT)
      --idp-provider=IDP-PROVIDER
                               The configured IDP provider. (env: SAML2ALIBABACLOUD_IDP_PROVIDER)
      --mfa=MFA                The name of the mfa. (env: SAML2ALIBABACLOUD_MFA)
  -s, --skip-verify            Skip verification of server certificate. (env: SAML2ALIBABACLOUD_SKIP_VERIFY)
      --url=URL                The URL of the SAML IDP server used to login. (env: SAML2ALIBABACLOUD_URL)
      --username=USERNAME      The username used to login. (env: SAML2ALIBABACLOUD_USERNAME)
      --password=PASSWORD      The password used to login. (env: SAML2ALIBABACLOUD_PASSWORD)
      --mfa-token=MFA-TOKEN    The current MFA token (supported in Keycloak, ADFS, GoogleApps). (env: SAML2ALIBABACLOUD_MFA_TOKEN)
      --role=ROLE              The ARN of the role to assume. (env: SAML2ALIBABACLOUD_ROLE)
      --urn=AlibabaCloudURN    The URN used by SAML when you login. (env: SAML2ALIBABACLOUD_URN)
      --skip-prompt            Skip prompting for parameters during login.
      --session-duration=SESSION-DURATION
                               The duration of your AlibabaCloud Session. (env: SAML2ALIBABACLOUD_SESSION_DURATION)
      --disable-keychain       Do not use keychain at all.
  -r, --region=REGION          AlibabaCloud region to use for API requests, e.g. cn-hangzhou (env: SAML2ALIBABACLOUD_REGION)

Commands:
  help [<command>...]
    Show help.


  configure [<flags>]
    Configure a new IDP account.

        --app-id=APP-ID            OneLogin app id required for SAML assertion. (env: ONELOGIN_APP_ID)
        --client-id=CLIENT-ID      OneLogin client id, used to generate API access token. (env: ONELOGIN_CLIENT_ID)
        --client-secret=CLIENT-SECRET
                                   OneLogin client secret, used to generate API access token. (env: ONELOGIN_CLIENT_SECRET)
        --subdomain=SUBDOMAIN      OneLogin subdomain of your company account. (env: ONELOGIN_SUBDOMAIN)
    -p, --profile=PROFILE          The AlibabaCloud CLI profile to save the temporary credentials. (env: SAML2ALIBABACLOUD_PROFILE)
        --resource-id=RESOURCE-ID  F5APM SAML resource ID of your company account. (env: SAML2ALIBABACLOUD_F5APM_RESOURCE_ID)
        --config=CONFIG            Path/filename of saml2alibabacloud config file (env: SAML2ALIBABACLOUD_CONFIGFILE)

  login [<flags>]
    Login to a SAML 2.0 IDP and convert the SAML assertion to an STS token.

    -p, --profile=PROFILE      The AlibabaCloud CLI profile to save the temporary credentials. (env: SAML2ALIBABACLOUD_PROFILE)
        --duo-mfa-option=DUO-MFA-OPTION
                               The MFA option you want to use to authenticate with
        --client-id=CLIENT-ID  OneLogin client id, used to generate API access token. (env: ONELOGIN_CLIENT_ID)
        --client-secret=CLIENT-SECRET
                               OneLogin client secret, used to generate API access token. (env: ONELOGIN_CLIENT_SECRET)
        --force                Refresh credentials even if not expired.

  exec [<flags>] [<command>...]
    Exec the supplied command with env vars from STS token.

    -p, --profile=PROFILE  The AlibabaCloud CLI profile to save the temporary credentials. (env: SAML2ALIBABACLOUD_PROFILE)
        --exec-profile=EXEC-PROFILE
                           The AlibabaCloud CLI profile to utilize for command execution. Useful to allow the AlibabaCloud cli to perform secondary role assumption. (env: SAML2ALIBABACLOUD_EXEC_PROFILE)

  console [<flags>]
    Console will open the AlibabaCloud console after logging in.

        --exec-profile=EXEC-PROFILE
                           The AlibabaCloud CLI profile to utilize for console execution. (env: SAML2ALIBABACLOUD_EXEC_PROFILE)
    -p, --profile=PROFILE  The AlibabaCloud CLI profile to save the temporary credentials. (env: SAML2ALIBABACLOUD_PROFILE)
        --force            Refresh credentials even if not expired.

  list-roles
    List available role ARNs.


  script [<flags>]
    Emit a script that will export environment variables.

    -p, --profile=PROFILE  The AlibabaCloud CLI profile to save the temporary credentials. (env: SAML2ALIBABACLOUD_PROFILE)
        --shell=bash       Type of shell environment. Options include: bash, powershell, fish


```


### `saml2alibabacloud script`

If the `script` sub-command is called, `saml2alibabacloud` will output the following temporary security credentials:
```
export ALIBABA_CLOUD_ACCESS_KEY_ID="STS.NTh..............8xN"
export ALIBABA_CLOUD_ACCESS_KEY_SECRET="CYc................T5M"
export ALIBABA_CLOUD_SESSION_TOKEN=""
export ALIBABA_CLOUD_SECURITY_TOKEN="CAI................xxP"
export ALICLOUD_ACCESS_KEY="STS.NTh................8xN"
export ALICLOUD_SECRET_KEY="CYc................T5M"
export ALICLOUD_SECURITY_TOKEN="CAI................xxP"
export ALICLOUD_PROFILE="saml"
export SAML2ALIBABA_CLOUD_PROFILE="saml"
```

Powershell, and fish shells are supported as well.

If you use `eval $(saml2alibabacloud script)` frequently, you may want to create a alias for it:

zsh:
```
alias s2a="function(){eval $( $(command saml2alibabacloud) script --shell=bash --profile=$@);}"
```

bash:
```
function s2a { eval $( $(which saml2alibabacloud) script --shell=bash --profile=$@); }
```

### `saml2alibabacloud exec`

If the `exec` sub-command is called, `saml2alibabacloud` will execute the command given as an argument:
By default saml2alibabacloud will execute the command with temp credentials generated via `saml2alibabacloud login`.

The `--exec-profile` flag allows for a command to execute using an AlibabaCloud CLI profile which may have chained
"assume role" actions. (via 'source_profile' in ~/.aliyun/config.json)

```
options:
--exec-profile           Execute the given command utilizing a specific profile from your ~/.aliyun/config.json file
```

### Configuring IDP Accounts

This is the *new* way of adding IDP provider accounts, it enables you to have named accounts with whatever settings you like and supports having one *default* account which is used if you omit the account flag. This replaces the --provider flag and old configuration file in 1.x.

To add a default IdP account to saml2alibabacloud just run the following command and follow the prompts.

```
$ saml2alibabacloud configure
? Please choose a provider: Ping
? AlibabaCloud CLI Profile myaccount

? URL https://example.com
? Username me@example.com

? Password
No password supplied

account {
  URL: https://example.com
  Username: me@example.com
  Provider: Ping
  MFA: Auto
  SkipVerify: false
  AlibabaCloudURN: urn:alibaba:cloudcomputing
  SessionDuration: 3600
  Profile: myaccount
  Region: us-east-1
}

Configuration saved for IDP account: default
```

Then to login using this account.

```
saml2alibabacloud login
```

You can also add named accounts, below is an example where I am setting up an account under the `wolfeidau` alias, again just follow the prompts.

```
saml2alibabacloud configure -a wolfeidau
```

You can also configure the account alias without prompts.

```
saml2alibabacloud configure -a wolfeidau --idp-provider KeyCloak --username mark@wolfe.id.au -r cn-north-1  \
  --url https://keycloak.wolfe.id.au/auth/realms/master/protocol/saml/clients/alibabacloud --skip-prompt
```

Then your ready to use saml2alibabacloud.

## Example

Log into a service (without MFA).

```
$ saml2alibabacloud login
Using IDP Account default to access Ping https://id.example.com
To use saved password just hit enter.
Username [mark.wolfe@example.com]:
Password: ************

Authenticating as mark.wolfe@example.com ...
Selected role: acs:ram::123123123123:role/Ali-Admin-CloudOPSNonProd
Requesting AlibabaCloud credentials using SAML assertion
Saving credentials
Logged in as: acs:ram::123123123123:role/Ali-Admin-CloudOPSNonProd/wolfeidau@example.com

Your new access key pair has been stored in the AlibabaCloud CLI configuration
To use this credential, call the AlibabaCloud CLI with the --profile option (e.g. aliyun --profile saml sts GetCallerIdentity --region=cn-hangzhou).
```

Log into a service (with MFA).

```
$ saml2alibabacloud login
Using IDP Account default to access Ping https://id.example.com
To use saved password just hit enter.
Username [mark.wolfe@example.com]:
Password: ************

Authenticating as mark.wolfe@example.com ...
Enter passcode: 123456

Selected role: acs:ram::123123123123:role/Ali-Admin-CloudOPSNonProd
Requesting AlibabaCloud credentials using SAML assertion
Saving credentials
Logged in as: acs:ram::123123123123:role/Ali-Admin-CloudOPSNonProd/wolfeidau@example.com

Your new access key pair has been stored in the AlibabaCloud CLI configuration
To use this credential, call the AlibabaCloud CLI with the --profile option (e.g. aliyun --profile saml sts GetCallerIdentity --region=cn-hangzhou).
```

## Advanced Configuration

Configuring multiple accounts with custom role and profile in `~/.aliyun/config.json` with goal being isolation between infra code when deploying to these environments. This setup assumes you're using separate roles and probably AlibabaCloud accounts for `dev` and `test` and is designed to help operations staff avoid accidentally deploying to the wrong AlibabaCloud account in complex environments. Note that this method configures SAML authentication to each AlibabaCloud account directly (in this case different AlibabaCloud accounts). In the example below, separate authentication values are configured for AlibabaCloud accounts 'profile=customer-dev/Account=was 121234567890' and 'profile=customer-test/Account=121234567891'

### Dev Account Setup

To setup the dev account run the following and enter URL, username and password, and assign a standard role to be automatically selected on login.

```
saml2alibabacloud configure -a customer-dev --role=acs:ram::121234567890:role/customer-admin-role -p customer-dev
```

This will result in the following configuration in `~/.saml2alibabacloud`.

```
[customer-dev]
url                              = https://id.customer.cloud
username                         = mark@wolfe.id.au
provider                         = Ping
mfa                              = Auto
skip_verify                      = false
timeout                          = 0
alibabacloud_urn                 = urn:alibaba:cloudcomputing
alibabacloud_session_duration    = 28800
alibabacloud_profile             = customer-dev
role_arn                         = acs:ram::121234567890:role/customer-admin-role
region                           = cn-hangzhou
```

To use this you will need to export `ALIBABACLOUD_DEFAULT_PROFILE=customer-dev` environment variable to target `dev`.

### Test Account Setup

To setup the test account run the following and enter URL, username and password.

```
saml2alibabacloud configure -a customer-test --role=acs:ram::121234567891:role/customer-admin-role -p customer-test
```

This results in the following configuration in `~/.saml2alibabacloud`.

```
[customer-test]
url                              = https://id.customer.cloud
username                         = mark@wolfe.id.au
provider                         = Ping
mfa                              = Auto
skip_verify                      = false
timeout                          = 0
alibabacloud_urn                 = urn:alibaba:cloudcomputing
alibabacloud_session_duration    = 28800
alibabacloud_profile             = customer-test
role_arn                         = acs:ram::121234567891:role/customer-admin-role
region                           = cn-hangzhou
```

To use this you will need to export `ALIBABACLOUD_DEFAULT_PROFILE=customer-test` environment variable to target `test`.


## Advanced Configuration - additional parameters
There are few additional parameters allowing to customise saml2alibabacloud configuration.
Use following parameters in `~/.saml2alibabacloud` file:
- `http_attempts_count` - configures the number of attempts to send http requests in order to authorise with saml provider. Defaults to 1
- `http_retry_delay` - configures the duration (in seconds) of timeout between attempts to send http requests to saml provider. Defaults to 1
- `region` - configures which region endpoints to use. Defaults to `cn-hangzhou`

Example: typical configuration with such parameters would look like follows:
```
[default]
url                              = https://id.customer.cloud
username                         = test@example.com
provider                         = Ping
mfa                              = Auto
skip_verify                      = false
timeout                          = 0
alibabacloud_urn                 = urn:alibaba:cloudcomputing
alibabacloud_session_duration    = 28800
alibabacloud_profile             = customer-dev
role_arn                         = acs:ram::121234567890:role/customer-admin-role
http_attempts_count              = 3
http_retry_delay                 = 1
region                           = cn-hangzhou
```
## Building

To build this software on osx clone to the repo to `$GOPATH/src/github.com/daxingplay/saml2alicloud` and ensure you have `$GOPATH/bin` in your `$PATH`.

```
make mod
```

Install the binary to `$GOPATH/bin`.

```
make install
```

Then to test the software just run.

```
make test
```

## Environment vars

The exec sub command will export the following environment variables.

* ALICLOUD_ACCESS_KEY
* ALICLOUD_SECRET_KEY
* ALICLOUD_SECURITY_TOKEN
* ALICLOUD_ASSUME_ROLE_SESSION_NAME
* ALICLOUD_PROFILE

Note: That profile environment variables enable you to use `exec` with a script or command which requires an explicit profile.

## Provider Specific Documentation

* [Azure Active Directory](./doc/provider/aad)
* [JumpCloud](./doc/provider/jumpcloud)

# Dependencies

This tool would not be possible without some great opensource libraries.

* [saml2aws](https://github.com/Versent/saml2aws) saml2aws
* [goquery](https://github.com/PuerkitoBio/goquery) html querying
* [etree](https://github.com/beevik/etree) xpath selector
* [kingpin](https://github.com/alecthomas/kingpin) command line flags
* [go-ini](https://github.com/go-ini/ini) INI file parser
* [go-ntlmssp](https://github.com/Azure/go-ntlmssp) NTLM/Negotiate authentication

# Releasing

Install `github-release`.

```
go get github.com/buildkite/github-release
```

To release run.

```
make release
```

# Debugging Issues with IDPs

There are two levels of debugging, first emits debug information and the URL / Method / Status line of requests.

```
saml2alibabacloud login --verbose
```

The second emits the content of requests and responses, this includes authentication related information so don't copy and paste it into chat or tickets!

```
DUMP_CONTENT=true saml2alibabacloud login --verbose
```

# License

This code is Copyright (c) 2018 [Versent](http://versent.com.au) and released under the MIT license. All rights not explicitly granted in the MIT license are reserved. See the included LICENSE.md file for more details.

