package encrypt

import (
	"crypto/md5"
	"fmt"

	"github.com/mangohow/cloud-ide/pkg/utils"
)

func Md5String(str string) string {
	return fmt.Sprintf("%x", md5.Sum(utils.String2Bytes(str)))
}
