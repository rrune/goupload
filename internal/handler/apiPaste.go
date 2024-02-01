package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/internal/database"
	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
)

func (h handler) Paste(c *fiber.Ctx) error {
	paste, err := h.DB.Pastes.GetPasteByShort(c.Params("short"))
	if util.CheckWLogs(err) {
		c.SendStatus(500)
	}

	return c.Render("paste", fiber.Map{
		"Paste": paste,
		"Url":   h.Url,
	})
}

func (h handler) GetPasteRaw(c *fiber.Ctx) error {
	paste, err := h.DB.Pastes.GetPasteByShort(c.Params("short"))
	if util.CheckWLogs(err) {
		c.SendStatus(500)
	}

	return c.SendString(paste.Text)
}

func (h handler) HandlePaste(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	username := claims["username"].(string)
	restrictedPerms := claims["restricted"].(bool)
	onetime := claims["onetime"].(bool)

	restricted := c.FormValue("restricted") == "restricted"
	if restricted && !restrictedPerms { // if restricted is requested but not permitted
		return c.SendStatus(403)
	}

	text := c.FormValue("text")
	if text == "" {
		return c.Render("response", fiber.Map{
			"Text":        "Textarea is empty",
			"Destination": "/",
		})
	}

	short, err := h.DB.Pastes.AddNewPaste(models.Paste{
		Text:       text,
		Author:     username,
		Timestamp:  time.Now(),
		Ip:         c.IP(),
		Restricted: restricted,
	})
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	if onetime {
		h.DB.Users.RemoveUserByUsername(username)
		c.ClearCookie("JWT")
	}

	return c.Render("response", fiber.Map{
		"Text":        "Pasted successfully",
		"Link":        h.Url + short,
		"Destination": "/",
	})
}

func (h handler) HandleRemovePaste(c *fiber.Ctx) error {
	short := c.Query("short", "")
	exist, err := database.CheckIfShortExists(h.DB.Pastes.DB, short)
	if util.CheckWLogs(err) || !exist {
		return c.SendStatus(400)
	}

	err = h.DB.Pastes.RemovePasteByShort(short)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	return c.Redirect("/dashboard")
}
