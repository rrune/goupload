package database

import (
	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(c models.Config) Database {
	dsn := c.DBUsername + ":" + c.DBPassword + "@" + c.DBAddress + "?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	util.CheckPanic(err)
	d := Database{
		Users:  UserDB{db},
		Shorts: ShortDB{db},
		Files:  FilesDB{db, c.ShortLength},
		Pastes: PastesDB{db, c.ShortLength},
	}
	return d
}

type Database struct {
	Users  UserDB
	Shorts ShortDB
	Files  FilesDB
	Pastes PastesDB
}
