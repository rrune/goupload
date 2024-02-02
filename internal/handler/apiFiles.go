package handler

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/internal/database"
	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
)

func (h handler) Download(c *fiber.Ctx) error {
	file, err := h.DB.Files.GetFileByShort(c.Params("short"))
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	return c.Download("./data/uploads/" + file.Filename)
}

func (h handler) Upload(c *fiber.Ctx, username string, blindPerms bool, restrictedPerms bool, onetime bool, callback func(filename string, short string, blind bool) error) error {
	blind := c.FormValue("blind") == "blind"
	if blind && !blindPerms { // if blind is requested but not permitted
		return c.SendStatus(403)
	}

	restricted := c.FormValue("restricted") == "restricted"
	if restricted && !restrictedPerms { // if restricted is requested but not permitted
		return c.SendStatus(403)
	}

	file, err := c.FormFile("file")
	if util.CheckWLogs(err) {
		return c.SendStatus(400)
	}

	var short string
	switch blind {
	case true:

		file.Filename, err = util.EnsureUniqueFilenames("./data/blind/", file.Filename)
		if util.CheckWLogs(err) {
			return c.SendStatus(500)
		}

		err = c.SaveFile(file, "./data/blind/"+file.Filename)

	case false:

		file.Filename, err = util.EnsureUniqueFilenames("./data/uploads/", file.Filename)
		if util.CheckWLogs(err) {
			return c.SendStatus(500)
		}

		short, err = h.DB.Files.AddNewFile(models.File{
			Filename:   file.Filename,
			Author:     username,
			Timestamp:  time.Now(),
			Ip:         c.IP(),
			Restricted: restricted,
		})
		if err == nil {
			err = c.SaveFile(file, "./data/uploads/"+file.Filename)
		}

	}
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	if onetime {
		h.DB.Users.RemoveUserByUsername(username)
		c.ClearCookie("JWT")
	}

	return callback(file.Filename, short, blind)
}

func (h handler) HandleUploadWeb(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	username := claims["username"].(string)
	blind := claims["blind"].(bool)
	restricted := claims["restricted"].(bool)
	onetime := claims["onetime"].(bool)

	return h.Upload(c, username, blind, restricted, onetime, func(filename, short string, blind bool) error {
		if blind {
			return c.Render("response", fiber.Map{
				"Text":        "Uploaded " + filename,
				"Destination": "/",
			})
		}
		return c.Render("response", fiber.Map{
			"Text":        "Uploaded " + filename,
			"Link":        h.Url + short,
			"Destination": "/",
		})
	})
}

// simpler upload to use with curl
func (h handler) HandleUploadSimple(c *fiber.Ctx) error {
	// uses simple auth, there are no other permission checks prior to this one
	authBase64 := strings.Split(c.Get("Authorization"), " ")

	// The slice's length needs to be 2, otherwise the header was empty
	if len(authBase64) < 2 {
		c.SendStatus(403)
	}
	// [0] of slice is "Basic"
	authByte, err := base64.StdEncoding.DecodeString(authBase64[1])
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	//basic auth is formated username:password
	auth := strings.Split(string(authByte), ":")
	if len(auth) != 2 {
		c.SendStatus(403)
	}
	username := auth[0]
	password := auth[1]

	correct, user, err := h.ValidateCredentials(username, password)
	if util.CheckWLogs(err) {
		c.SendStatus(500)
	}
	if !correct {
		c.SendStatus(403)
	}

	return h.Upload(c, username, user.Blind, user.Restricted, user.Onetime, func(filename, short string, blind bool) error {
		if blind {
			return c.SendString("Uploaded " + filename)
		}
		return c.SendString("Uploaded " + filename + "\n" + h.Url + short)
	})
}

func (h handler) HandleRemoveFile(c *fiber.Ctx) error {
	short := c.Params("short", "")

	exist, err := database.CheckIfShortExists(h.DB.Files.DB, short)
	if util.CheckWLogs(err) || !exist {
		return c.SendStatus(400)
	}

	file, err := h.DB.Files.GetFileByShort(short)
	if err == nil {

		err = os.Remove("./data/uploads/" + file.Filename)
		if err == nil {

			err = h.DB.Files.RemoveFileByShort(short)
			if err == nil {
				return c.Redirect("/dashboard")
			}
		}
	}
	util.CheckWLogs(err)
	return c.SendStatus(500)
}

func (h handler) HandleMoveToBlind(c *fiber.Ctx) error {
	short := c.Params("short", "")
	if short == "" {
		return c.SendStatus(400)
	}
	file, err := h.DB.Files.GetFileByShort(short)
	if err == nil {
		err = os.Rename("./data/uploads/"+file.Filename, "./data/blind/"+file.Filename)
		if err == nil {
			err = h.DB.Files.RemoveFileByShort(short)
		}
	}
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	return c.Redirect("/dashboard")
}

func (h handler) HandleDetails(c *fiber.Ctx) error {
	short := c.Params("short", "")
	if short == "" {
		return c.SendStatus(400)
	}
	file, err := h.DB.Files.GetFileByShort(short)
	var info os.FileInfo
	if err == nil {
		info, err = os.Stat("./data/uploads/" + file.Filename)
	}
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	infostrings := []string{
		info.Name(),
		strconv.FormatFloat(float64(info.Size())/1000, 'f', -1, 64) + " kb",
		strconv.Itoa(file.Downloads) + " Downloads",
		file.Ip,
	}
	return c.Render("response", fiber.Map{
		"Strings":     infostrings,
		"Destination": "/dashboard",
	})
}
