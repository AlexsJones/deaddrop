/*********************************************************************************
*     File Name           :     utils/hash.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-15 21:29]
*     Last Modified       :     [2015-10-15 21:34]
*     Description         :      
**********************************************************************************/

package utils

import (
  "crypto/md5"
  "fmt"
  "encoding/hex"
)

func Hash(input string) string {
  h := md5.New()
  h.Write([]byte(input))
  return fmt.Sprintf(hex.EncodeToString(h.Sum(nil)))
}
