/*********************************************************************************
*     File Name           :     utils/crypto.go
*     Created By          :     tibbar
*     Creation Date       :     [2015-10-22 10:35]
*     Last Modified       :     [2015-10-22 11:33]
*     Description         :      
**********************************************************************************/

package utils
import (
  "crypto/aes"
  "crypto/cipher"
  "strings"
  "os"
  "bytes"
  "log"
  "io"
  "crypto/rand"
  "io/ioutil"
)

func DecryptContent(filePath string, key []byte) string {

  raw, err := ioutil.ReadFile(filePath)
  if err != nil {
    log.Println("Error in reading file for decryption")
    return ""
  }  

  block, err := aes.NewCipher(key)
  if err != nil {
    log.Println(err.Error())
    return ""
  }  

  ciphertext := []byte(raw)
  
  if len(ciphertext) < aes.BlockSize {
    log.Println("Error in DecryptContent size")
    return ""
  }
  
  iv := ciphertext[:aes.BlockSize]

  ciphertext = ciphertext[aes.BlockSize:]

  stream := cipher.NewCFBDecrypter(block,iv)

  stream.XORKeyStream(ciphertext,ciphertext)

  if strings.Contains(filePath,".encrypt") == false {
    log.Println("The file did not have a .encrypt extension")
    return ""
  }

  outputPath := strings.Split(filePath,".encrypt")[0]

  f, err := os.Create(outputPath)

  if err != nil {
    panic(err.Error())
  }
  _, err = io.Copy(f, bytes.NewReader(ciphertext))
  if err != nil {
    panic(err.Error())
  }

  log.Println("Deleting " + filePath)
  os.Remove(filePath)

  return outputPath
}

func EncryptContent(filePath string, key []byte) string {

  raw, err := ioutil.ReadFile(filePath)
  if err != nil {
    panic(err)
  }  

  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }  

  ciphertext := make([]byte, aes.BlockSize+len(raw))
  iv := ciphertext[:aes.BlockSize]
  if _, err := io.ReadFull(rand.Reader,iv); err != nil {
    panic(err)
  }

  stream := cipher.NewCFBEncrypter(block,iv)
  stream.XORKeyStream(ciphertext[aes.BlockSize:], raw)

  f, err := os.Create(filePath + ".encrypt")

  if err != nil {
    panic(err.Error())
  }

  _, err = io.Copy(f, bytes.NewReader(ciphertext))
  if err != nil {
    panic(err.Error())
  }

  log.Println("Deleting " + filePath)
  os.Remove(filePath)

  return filePath + ".encrypt"
}
