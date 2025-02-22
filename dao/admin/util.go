package adminDAO

import (
	"crypto/sha256"
	"fmt"
)

func genSaltPassword(salt string, password string) string {
	sh1 := sha256.New()
	sh1.Write([]byte(salt))
	sh2 := sha256.New()
	sh2.Write([]byte(fmt.Sprintf("%x", sh1.Sum(nil)) + password))
	return fmt.Sprintf("%x", sh2.Sum(nil))
}
