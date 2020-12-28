# EAA Akamai IDP support for saml2alibabacloud
This code supports to authenticate user from cli with saml to Akamai SAML IDP without browser

# Requirements
* Need Go 1.12
* Need saml2alibabacloud code from github link https://github.com/aliyun/saml2alibabacloud

# Building the Code
* Install Go 1.12
* Set GOPATH
* clone code from github link to $GOPATH/src/github.com/aliyun/saml2alibabacloud
* copy akamai.go to aliyun/saml2alibabacloud/pkg/providers/akamai/
* Merge code from saml2alibabacloud.go to support Akamai config.
* Ensure $GOPATH/bin in your $PATH
* make deps
* make install
* Binary will be present in GOPATH/bin

# Configuring the SAML IDP
* create Akamai EAA IDP
* Create a saml saas app
* Add Attribute as mentioned below in example.

"attrmap": [
     {
          "fmt": "uri_reference",
          "name": "https://www.aliyun.com/SAML-Role/Attributes/Role",
          "src": "",
          "val": "acs:ram::432929478872:saml-provider/AkamaiIDP,acs:ram::432929478872:role/AkamaiIDProle"
     },
     {
          "fmt": "basic",
          "name": "https://www.aliyun.com/SAML-Role/Attributes/RoleSessionName",
          "val": "punit@qadomain.com"
     },
     {
          "fmt": "basic",
          "name": "https://www.aliyun.com/SAML-Role/Attributes/SessionDuration",
          "val": "1200"
     }
]

# Using the saml2alibabacloud
* Configure IDP account run command -  saml2alibabacloud configure.
* Add url as https://<EAAIDP>/?app=<SAAShostname> Eg: https://samlidp.example.com/?app=signin.aliyun.com
* To login using saml2alibabacloud run command - saml2alibabacloud login
