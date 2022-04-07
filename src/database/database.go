package database

import (
	"github.com/rrune/goupload/models"
	"github.com/rrune/goupload/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(c models.Config) Database {
	dsn := c.Username + ":" + c.Password + "@" + c.Address + "?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	util.CheckPanic(err)
	d := Database{
		Users: UserDB{db},
		Files: FilesDB{db, c.ShortLength},
	}
	return d
}

type Database struct {
	Users UserDB
	Files FilesDB
}
