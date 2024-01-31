package handler

import (
	"slices"

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
		"Username":   claims["username"].(string),
		"Root":       claims["root"].(bool),
		"Blind":      claims["blind"].(bool),
		"Restricted": claims["restricted"].(bool),
	})
}

func (t template) Dashboard(c *fiber.Ctx) error {
	users, err := t.DB.Users.GetAllUsers()
	if util.Check(err) {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	files, err := t.DB.Files.GetAllFiles()
	if util.Check(err) {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// reverse the slice to show newest first
	slices.Reverse[[]models.File](files)
	return c.Render("dashboard", fiber.Map{
		"Users": users,
		"Files": files,
		"Url":   t.Url,
	})
}

func (t template) Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Msg": c.Query("msg", ""),
	})
}
