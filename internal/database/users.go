package database

import (
	"errors"

	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func (d UserDB) GetAllUsers() (u []models.User, err error) {
	err = d.DB.Table("Users").Find(&u).Error
	return
}

func (d UserDB) GetUserByUsername(username string) (u models.User, err error) {
	err = d.DB.Table("Users").First(&u, "Username = ?", username).Error
	return
}

func (d UserDB) CreateUser(user *models.User) (err error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if util.Check(err) {
		return
	}
	user.Password = hashedPassword

	_, err2 := d.GetUserByUsername(user.Username)
	if !util.Check(err2) {
		err = errors.New("username already exists")
		return
	}

	err = d.DB.Table("Users").Create(&user).Error
	return
}

func (d UserDB) RemoveUserByUsername(username string) (err error) {
	allUsers, err := d.GetAllUsers()
	if util.Check(err) {
		return
	}
	if len(allUsers) == 1 {
		err = errors.New("cannot remove only user")
		return
	}

	err = d.DB.Table("Users").Where("Username = ?", username).Delete(&models.User{}).Error
	return
}

func (d UserDB) ChangePassword(username string, password string) (err error) {
	hashedPassword, err := util.HashPassword(password)
	if util.Check(err) {
		return
	}
	err = d.DB.Table("Users").Where("Username = ?", username).Update("Password", hashedPassword).Error
	return
}

func (d UserDB) ChangePerms(user models.User) (err error) {
	err = d.DB.Table("Users").Where("Username = ?", user.Username).Update("Root", user.Root).Update("Blind", user.Blind).Update("Restricted", user.Restricted).Update("Onetime", user.Onetime).Error
	return
}
