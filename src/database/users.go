package database

import (
	"errors"

	"github.com/rrune/goupload/models"
	"github.com/rrune/goupload/util"
	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func (d UserDB) GetAllUsers() (u []models.User, err error) {
	err = d.DB.Table("uploaderUsers").Find(&u).Error
	return
}

func (d UserDB) GetUserByUsername(username string) (u models.User, err error) {
	err = d.DB.Table("uploaderUsers").First(&u, "username = ?", username).Error
	return
}

func (d UserDB) CreateUser(user *models.User) (err error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if util.Check(err) {
		return
	}
	user.Password = hashedPassword

	_, err2 := d.GetUserByUsername(user.Username)
	if err2 == nil {
		err = errors.New("Username already exists")
		return
	}

	err = d.DB.Table("uploaderUsers").Create(&user).Error
	return
}

func (d UserDB) RemoveUserByUsername(username string) (err error) {
	allUsers, err := d.GetAllUsers()
	if util.Check(err) {
		return
	}
	if len(allUsers) == 1 {
		err = errors.New("Cannot remove only user")
		return
	}

	err = d.DB.Table("uploaderUsers").Where("username = ?", username).Delete(&models.User{}).Error
	return
}

func (d UserDB) ChangePassword(username string, password string) (err error) {
	hashedPassword, err := util.HashPassword(password)
	if util.Check(err) {
		return
	}
	err = d.DB.Table("uploaderUsers").Where("username = ?", username).Update("password", hashedPassword).Error
	return
}
