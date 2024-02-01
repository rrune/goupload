package handler

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/internal/database"
	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
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
	if util.CheckWLogs(err) {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	files, err := t.DB.Files.GetAllFiles()
	if util.CheckWLogs(err) {
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

func (t template) AddUser(c *fiber.Ctx) error {
	return c.Render("createUser", fiber.Map{})
}

func (t template) ChangePassword(c *fiber.Ctx) error {
	username := c.Query("username", "")
	if username == "" {
		return c.SendStatus(400)
	}
	return c.Render("changePassword", fiber.Map{
		"Username": username,
	})
}

func (t template) ChangePerms(c *fiber.Ctx) error {
	username := c.Query("username", "")
	if username == "" {
		return c.SendStatus(400)
	}
	user, err := t.DB.Users.GetUserByUsername(username)
	if util.CheckWLogs(err) {
		return c.Render("response", fiber.Map{
			"Text":        "Database error",
			"Destination": "/dashboard",
		})
	}
	user.Password = ""
	return c.Render("changePerms", fiber.Map{
		"Username": username,
		"User":     user,
	})
}

func (t template) Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Msg": c.Query("msg", ""),
	})
}
