package commands

import (
	"github.com/aliyun/saml2alibabacloud/helper/credentials"
	"github.com/aliyun/saml2alibabacloud/helper/linuxkeyring"
)

func init() {
	if keyringHelper, err := linuxkeyring.NewKeyringHelper(); err == nil {
		credentials.CurrentHelper = keyringHelper
	}
}
