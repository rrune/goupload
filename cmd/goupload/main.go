package main

import (
	"log"
	"os"

	"github.com/rrune/goupload/internal/database"
	"github.com/rrune/goupload/internal/handler"
	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
	"gopkg.in/yaml.v2"
)

func main() {
	f, err := os.OpenFile("./data/goupload.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(2 | 3)

	var config models.Config
	ymlData, err := os.ReadFile("./config/config.yml")
	util.CheckPanic(err)
	err = yaml.Unmarshal(ymlData, &config)
	util.CheckPanic(err)

	db := database.New(config)

	// create new user if no user exists
	users, err := db.Users.GetAllUsers()
	util.CheckPanic(err)
	if len(users) == 0 {
		db.Users.CreateUser(&models.User{
			Username:   config.Username,
			Password:   config.Password,
			Root:       true,
			Blind:      true,
			Onetime:    false,
			Restricted: true,
		})
	}

	handler.Start(config.Port, config.JWTKey, config.Url, config.UploadLimit, db)
}
