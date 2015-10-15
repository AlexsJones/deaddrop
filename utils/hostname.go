/*********************************************************************************
*     File Name           :     utils/hostname.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-15 08:30]
*     Last Modified       :     [2015-10-15 08:32]
*     Description         :      
**********************************************************************************/

package utils
import (
  "fmt"
  "os"
  "net"
)
func GetHostnameIpv4() string {

  host, _ := os.Hostname()
  addrs, _ := net.LookupIP(host)
  for _, addr := range addrs {
    if ipv4 := addr.To4(); ipv4 != nil {
      return fmt.Sprintf("%s",ipv4)
    }
  }
  return ""
}
