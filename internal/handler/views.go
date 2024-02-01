package handler

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/internal/database"
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
		return c.SendStatus(500)
	}
	files, err := t.DB.Files.GetAllFiles()
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	pastes, err := t.DB.Pastes.GetAllPastes()
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	for i, paste := range pastes {
		if len(paste.Text) > 150 {
			pastes[i].Text = paste.Text[:150] + "..."
		}
	}

	// sort the slice to show newest first
	sort.Slice(files, func(i, j int) bool {
		return files[i].Timestamp.After(files[j].Timestamp)
	})

	sort.Slice(pastes, func(i, j int) bool {
		return pastes[i].Timestamp.After(pastes[j].Timestamp)
	})

	return c.Render("dashboard", fiber.Map{
		"Users":  users,
		"Files":  files,
		"Pastes": pastes,
		"Url":    t.Url,
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
