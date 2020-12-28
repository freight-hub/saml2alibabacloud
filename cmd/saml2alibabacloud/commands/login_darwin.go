package commands

import (
	"github.com/aliyun/saml2alibabacloud/helper/credentials"
	"github.com/aliyun/saml2alibabacloud/helper/osxkeychain"
)

func init() {
	credentials.CurrentHelper = &osxkeychain.Osxkeychain{}
}
