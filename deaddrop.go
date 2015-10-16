package main

import (
  "bytes"
  _ "github.com/lib/pq"
  "flag"
  "fmt"
  "os"
  "io"
  "log"	
  "net/http"
  "github.com/gorilla/mux"
  "github.com/AlexsJones/deaddrop/utils"
  "io/ioutil"
  "html/template"
  "strings"
)

var port string 

var configuration utils.Configuration 


func hdeaddrop_upload(w http.ResponseWriter, r *http.Request) {
  
  if r.ContentLength > (1024 * 10) {

      http.Error(w,"request too large", http.StatusExpectationFailed)
      return
  }

  const _24K = (1 << 20) * 24  
  if err := r.ParseMultipartForm(_24K); nil != err {  
    return 
  } 

  file, header, err := r.FormFile("file") 
  defer file.Close()

  if err != nil {
    fmt.Fprintln(w, err)
    log.Println(err.Error())
    return
  }

  guid := utils.NewGuid()

  hashedGuid := utils.Hash(guid)

  filenameCipher := hashedGuid + "_" + header.Filename

  out, err := os.Create("uploads/" + filenameCipher)

  if err != nil {
    fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
    return
  }

  defer out.Close()

  _, err = io.Copy(out, file)
  if err != nil {
    fmt.Fprintln(w, err)
  }

  w.Write([]byte("Generating 1 time download code: "+ hashedGuid))

}

func hdeaddrop_fetch(w http.ResponseWriter, r *http.Request) {

  r.ParseForm()
  id := r.FormValue("id")

  dirname := "uploads"
  d, err := os.Open(dirname)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer d.Close()
  fi, err := d.Readdir(-1)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  for _, fi := range fi {
    if fi.Mode().IsRegular() {

      splitString := strings.Split(fi.Name(),"_") 

      if splitString[0] == id {

	file := "uploads/" + fi.Name()

	log.Println(file)

	http.ServeFile(w,r,file)

	/* Delete file */
	os.Remove(file)
      }
    }
  }
  w.Write([]byte("404"))
}

func hdeaddrop_home(w http.ResponseWriter, r *http.Request) {

  s1, _ := template.ParseFiles("tmpl/header.tmpl", 
  "tmpl/content.tmpl", "tmpl/footer.tmpl")

  fbody, err := ioutil.ReadFile("views/index.html")
  if err != nil {

  }
  var buffer bytes.Buffer
  s1.ExecuteTemplate(&buffer, "header", nil)
  s1.ExecuteTemplate(&buffer, "content", template.HTML(string(fbody)))
  s1.ExecuteTemplate(&buffer, "footer", nil)

  w.Write(buffer.Bytes())
}

func main() {

  os.RemoveAll("uploads")
  os.Mkdir("uploads",0777)

  if os.Getenv("DEADDROP_CONF")  != "" {
    configuration = utils.NewConfiguration(os.Getenv("DEADDROP_CONF"))
  }else {
    var confFlag = flag.String("conf","","Path to configuration file")

    flag.Parse()

    if *confFlag == "" {
      log.Fatal("Please provide a conf path -conf") 
      return
    }
    configuration = utils.NewConfiguration(*confFlag)
  }

  rtr := mux.NewRouter()

  rtr.HandleFunc("/deaddrop/fetch", hdeaddrop_fetch).Methods("POST")

  rtr.HandleFunc("/deaddrop/upload",hdeaddrop_upload).Methods("POST")

  rtr.HandleFunc("/",hdeaddrop_home).Methods("GET")

  http.Handle("/",rtr)

  publicfolder := http.FileServer(http.Dir("./public/"))

  http.Handle("/public/",http.StripPrefix("/public/",publicfolder))

  port = configuration.Json.Port
  if os.Getenv("PORT") != "" {
    port = os.Getenv("PORT")  
    log.Print("Using environmental variable for $PORT")
  }

  log.Println("Listening...")

  http.ListenAndServe(":" + port ,nil)
}
