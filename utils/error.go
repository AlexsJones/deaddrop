/*********************************************************************************
*     File Name           :     util/error.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-05 16:19]
*     Last Modified       :     [2015-11-18 12:59]
*     Description         :      
**********************************************************************************/

package utils

import "log"
func CheckErr(err error, msg string) {
  if err != nil {
    log.Fatalln(msg, err)
  }
}
