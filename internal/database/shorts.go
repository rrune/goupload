package database

import (
	"math/rand"

	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
	"gorm.io/gorm"
)

type ShortDB struct {
	DB *gorm.DB
}

func (d ShortDB) GetShort(short string) (r models.Short, err error) {
	err = d.DB.
		Table("Shorts").Where("Short = ?", short).
		First(&r).
		Error
	return
}

func (d ShortDB) SwitchRestrict(short string) (err error) {
	paste, err := d.GetShort(short)
	if util.Check(err) {
		return
	}
	err = d.DB.
		Table("Shorts").
		Where("Short = ?", short).
		Update("Restricted", !paste.Restricted).
		Error
	return
}

func (d ShortDB) UpdateDownloadCounter(short string, downloads int) (err error) {
	err = d.DB.
		Table("Shorts").
		Where("Short = ?", short).
		Update("Downloads", downloads+1).
		Error
	return
}

func CheckIfShortExists(db *gorm.DB, short string) (exists bool, err error) {
	result := []DB_Short{}
	r := db.
		Table("Shorts").
		Where("Short = ?", short).
		Limit(1).
		Find(&result)

	err = r.Error
	exists = r.RowsAffected > 0
	return
}

func getShort(db *gorm.DB, shortLength int) (short string, err error) {
	exists := true
	for exists {
		short = random(shortLength)
		exists, err = CheckIfShortExists(db, short)
		if util.Check(err) {
			return
		}
	}
	return
}

func random(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
