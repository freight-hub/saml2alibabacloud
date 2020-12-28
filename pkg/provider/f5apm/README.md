# F5 Access Policy Manager Provider

* https://www.f5.com/products/security/access-policy-manager

## Instructions

You'll need the SAML policy ID for the AlibabaCloud account.  Your admin should be able to 
provide this (or you'll briefly see it in a redirect when you click an application link)

```
https://<YOUR ORGS DOMAIN>/saml/idp/res?id=<SAML RESOURCE ID>
```

Example Config:

```
[default]
url                           = https://<YOUR ORGS DOMAIN>
username                      = <YOUR USERNAME>
provider                      = F5APM
mfa                           = Auto
skip_verify                   = false
timeout                       = 0
alibabacloud_urn              = urn:alibaba:cloudcomputing
alibabacloud_session_duration = 3600
alibabacloud_profile          = <AlibabaCloud CLI PROFILE NAME>
resource_id                   = <SAML RESOURCE ID>
role_arn                      = 
```

Where `resource_id` will be something like `/Common/example-alibabacloud-account`

## Features

* Automatic detection of MFA
* Automatic detection of MFA options (push, token)
