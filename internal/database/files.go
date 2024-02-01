package database

import (
	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
	"gorm.io/gorm"
)

type FilesDB struct {
	DB          *gorm.DB
	ShortLength int
}

func (d FilesDB) GetAllFiles() (f []models.File, err error) {
	err = d.DB.
		Table("Shorts").
		Select("Shorts.*, Files.Filename").
		Joins("INNER JOIN Files ON Shorts.Short = Files.Short").
		Find(&f).
		Error
	return
}

func (d FilesDB) GetFileByShort(short string) (f models.File, err error) {
	err = d.DB.
		Table("Shorts").Where("Shorts.Short = ?", short).
		Select("Shorts.*, Files.Filename").
		Joins("INNER JOIN Files ON Shorts.Short = Files.Short").
		Find(&f).
		Error
	return
}

func (d FilesDB) GetFileByName(filename string) (f models.File, err error) {
	err = d.DB.
		Table("Files").Where("Files.Filename = ?", filename).
		Select("Shorts.*, Files.Filename").
		Joins("INNER JOIN Shorts ON Files.Short = Shorts.Short").
		Find(&f).
		Error
	return
}

func (d FilesDB) AddNewFile(file models.File) (short string, err error) {
	short, err = getShort(d.DB, d.ShortLength)
	if util.Check(err) {
		return
	}

	dbFile := DB_File{
		Short:    short,
		Filename: file.Filename,
	}

	err = d.DB.
		Table("Files").
		Create(&dbFile).
		Error
	if util.Check(err) {
		return
	}

	dbShort := DB_Short{
		Short:      short,
		Type:       "file",
		Author:     file.Author,
		Timestamp:  file.Timestamp,
		Ip:         file.Ip,
		Restricted: file.Restricted,
	}
	err = d.DB.
		Table("Shorts").
		Create(&dbShort).
		Error
	return
}

func (d FilesDB) RemoveFileByShort(short string) (err error) {
	err = d.DB.
		Table("Files").
		Where("Short = ?", short).
		Delete(&models.File{}).
		Error
	if util.Check(err) {
		return
	}
	err = d.DB.
		Table("Shorts").
		Where("Short = ?", short).
		Delete(&models.File{}).
		Error
	return
}
