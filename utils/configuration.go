package utils

import (
  "os"
  "log"
  "encoding/json"
)
type Json struct {
  Port string
  DBConnectionString string
}

type Configuration struct {
  Json Json
}

func parseJson(configurationPath string) Json {

  var j Json
  conf,err := os.Open(configurationPath)
  if err != nil {
    log.Fatalf("opening configuration file",err.Error())
  }

  jsonParser := json.NewDecoder(conf)
  if err = jsonParser.Decode(&j); err != nil {
    log.Fatalf("parsing config file", err.Error())
  }
  return j
}

func NewConfiguration(configurationPath string) Configuration {

  c := Configuration {
    Json: parseJson(configurationPath),
  }
  return c
}
