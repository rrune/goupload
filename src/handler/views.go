package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/database"
	"github.com/rrune/goupload/models"
	"github.com/rrune/goupload/util"
)

type template struct {
	DB  database.Database
	Url string
}

func (t template) Index(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Render("index", fiber.Map{
		"Username": claims["username"].(string),
		"Root":     claims["root"].(bool),
		"Blind":    claims["blind"].(bool),
	})
}

func (t template) Manage(c *fiber.Ctx) error {
	users, err := t.DB.Users.GetAllUsers()
	if util.Check(err) {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	files, err := t.DB.Files.GetAllFiles()
	if util.Check(err) {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Render("manage", fiber.Map{
		"Users": users,
		"Files": files,
		"Url":   t.Url,
	})
}

func (t template) Filter(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	username := user.Username

	users, err := t.DB.Users.GetAllUsers()
	if util.Check(err) {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	files, err := t.DB.Files.GetAllFiles()
	if util.Check(err) {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var filtered []models.File
	for _, file := range files {
		if file.Author == username {
			filtered = append(filtered, file)
		}
	}
	return c.Render("filter", fiber.Map{
		"Users": users,
		"Files": filtered,
		"Url":   t.Url,
	})
}

func (t template) Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Msg": c.Query("msg", ""),
	})
}
