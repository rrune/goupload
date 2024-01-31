package main

import (
	"log"
	"os"

	"github.com/rrune/goupload/database"
	"github.com/rrune/goupload/handler"
	"github.com/rrune/goupload/models"
	"github.com/rrune/goupload/util"
	"gopkg.in/yaml.v2"
)

func main() {
	f, err := os.OpenFile("../data/goupload.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(2 | 3)

	var config models.Config
	ymlData, err := os.ReadFile("../data/config.yml")
	util.CheckPanic(err)
	err = yaml.Unmarshal(ymlData, &config)
	util.CheckPanic(err)

	db := database.New(config)

	handler.Start(config.Port, config.JWTKey, config.Url, config.UploadLimit, db)
}
