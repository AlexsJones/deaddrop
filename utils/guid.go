/*********************************************************************************
*     File Name           :     guid.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 15:06]
*     Last Modified       :     [2015-10-15 08:26]
*     Description         :      
**********************************************************************************/

package utils

import (
  "fmt"
  "log"
  "github.com/nu7hatch/gouuid"
)

func NewGuid() string {
  u, err := uuid.NewV4()
  if err != nil {
    log.Println("Unable to retrieve UUID")
    return ""
  }
  return fmt.Sprintf("%s",u)
}
