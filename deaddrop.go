package main

import (
  _ "github.com/lib/pq"
  "flag"
  "os"
  "log"	
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/AlexsJones/deaddrop/utils"
  "github.com/jinzhu/gorm"
)

var databaseConnect gorm.DB

type incoming_data struct {
  Data string
  Guid string
}

func hdeaddrop_post(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var i incoming_data
  err := decoder.Decode(&i)
  if err != nil {
    log.Println("Malformed post")
    return
  }
  log.Println(i.Data)

  guid := utils.NewGuid()

  /* Store data in db */
  i.Guid = guid

  databaseConnect.Create(&i)

  w.Write([]byte("Generating 1 time download link " + "http://" + utils.GetHostnameIpv4() + "/" + guid))
}

func hdeaddrop_get(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  id := params["id"]

  var result incoming_data

  databaseConnect.Where(&incoming_data{ Guid: id}).First(&result)

  databaseConnect.Where(&incoming_data{ Guid: id}).Delete(&incoming_data{})

  w.Write([]byte(result.Data))
}

func initialiseDatabase() {

  db, err := gorm.Open("postgres","host=localhost port=5432 user=anon password=anon dbname=development")

  utils.CheckErr(err, "postgres failed")

  db.DB()

  db.DB().Ping()

  db.SingularTable(true)

  db.CreateTable(&incoming_data{})

  databaseConnect = db

}
func main() {

  var configuration utils.Configuration 

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

  initialiseDatabase()

  rtr := mux.NewRouter()

  rtr.HandleFunc("/deaddrop/{id}", hdeaddrop_get).Methods("GET")

  rtr.HandleFunc("/deaddrop", hdeaddrop_post).Methods("POST")

  http.Handle("/",rtr)

 
  port := configuration.Json.Port
  if os.Getenv("PORT") != "" {
    port = os.Getenv("PORT")  
    log.Print("Using environmental variable for $PORT")
  }

 log.Println("Listening...")

  http.ListenAndServe(":" + port ,nil)
}
