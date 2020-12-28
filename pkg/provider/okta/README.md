# Okta provider

## Instructions

Retrieve the AlibabaCloud application URL from your Okta tenant. This will (may) look something like:

```
https://$YOUR_ORGANIZATION.okta.com/home/alibabacloud/$OKTA_APPLICATION_ID/$OKTA_OTHER_ID
```

The path segments `/home/alibabacloud` in the above URL may vary.

## Features

* Supports MFA (Okta Push, Okta TOTP, Duo, and Google Authenticator), when configured at *organization* or *application* level.