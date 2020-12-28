package commands

import (
	"github.com/aliyun/saml2alibabacloud/helper/credentials"
	"github.com/aliyun/saml2alibabacloud/helper/wincred"
)

func init() {
	credentials.CurrentHelper = &wincred.Wincred{}
}
