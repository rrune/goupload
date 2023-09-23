package database

import (
	"math/rand"
	"time"

	"github.com/rrune/goupload/models"
	"gorm.io/gorm"
)

type FilesDB struct {
	DB          *gorm.DB
	ShortLength int
}

func (d FilesDB) GetAllFiles() (u []models.File, err error) {
	err = d.DB.Table("uploadedFiles").Find(&u).Error
	return
}

func (d FilesDB) GetFileByShort(short string) (f models.File, err error) {
	err = d.DB.Table("uploadedFiles").First(&f, "short = ?", short).Error
	return
}

func (d FilesDB) AddNewFile(file models.File) (short string, err error) {
	short = d.getShort()
	file.Short = short

	err = d.DB.Table("uploadedFiles").Create(&file).Error
	return
}

func (d FilesDB) RemoveFileByShort(short string) (err error) {
	err = d.DB.Table("uploadedFiles").Where("short = ?", short).Delete(&models.File{}).Error
	return
}

func (d FilesDB) SwitchRestrict(short string) (err error) {
	file, err := d.GetFileByShort(short)
	if err != nil {
		return
	}
	err = d.DB.Table("uploadedFiles").Where("short = ?", short).Update("restricted", !file.Restricted).Error
	return
}

func (d FilesDB) getShort() string {
	var random string
	unique := false
	for !unique {
		random = d.random(d.ShortLength)
		_, err := d.GetFileByShort(random)
		if err != nil {
			unique = true
		}
	}
	return random
}

func (d FilesDB) random(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
