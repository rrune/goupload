package database

import (
	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
	"gorm.io/gorm"
)

type PastesDB struct {
	DB          *gorm.DB
	ShortLength int
}

func (d PastesDB) GetAllPastes() (p []models.Paste, err error) {
	err = d.DB.
		Table("Shorts").
		Select("Shorts.*, Pastes.Text").
		Joins("INNER JOIN Pastes ON Shorts.Short = Pastes.Short").
		Find(&p).
		Error
	return
}

func (d PastesDB) GetPasteByShort(short string) (p models.Paste, err error) {
	err = d.DB.
		Table("Shorts").Where("Shorts.Short = ?", short).
		Select("Shorts.*, Pastes.Text").
		Joins("INNER JOIN Pastes ON Shorts.Short = Pastes.Short").
		Find(&p).
		Error
	return
}

func (d PastesDB) AddNewPaste(paste models.Paste) (short string, err error) {
	short, err = getShort(d.DB, d.ShortLength)
	if util.Check(err) {
		return
	}

	dbPaste := DB_Paste{
		Short: short,
		Text:  paste.Text,
	}

	err = d.DB.
		Table("Pastes").
		Create(&dbPaste).
		Error
	if util.Check(err) {
		return
	}

	dbShort := DB_Short{
		Short:      short,
		Type:       "paste",
		Author:     paste.Author,
		Timestamp:  paste.Timestamp,
		Ip:         paste.Ip,
		Restricted: paste.Restricted,
	}
	err = d.DB.
		Table("Shorts").
		Create(&dbShort).
		Error
	return
}

func (d PastesDB) RemovePasteByShort(short string) (err error) {
	err = d.DB.
		Table("Pastes").
		Where("Short = ?", short).
		Delete(&models.Paste{}).
		Error
	if util.Check(err) {
		return
	}
	err = d.DB.
		Table("Shorts").
		Where("Short = ?", short).
		Delete(&models.Paste{}).
		Error
	return
}

func (d PastesDB) SwitchRestrict(short string) (err error) {
	paste, err := d.GetPasteByShort(short)
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

func (d PastesDB) UpdateDownloadCounter(short string, downloads int) (err error) {
	err = d.DB.
		Table("Shorts").
		Where("Short = ?", short).
		Update("Downloads", downloads+1).
		Error
	return
}
