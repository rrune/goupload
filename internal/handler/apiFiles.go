package handler

import (
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
	"gorm.io/gorm"
)

// edit
func (h handler) HandleDownload(c *fiber.Ctx) error {
	file, err := h.DB.Files.GetFileByShort(c.Params("short"))
	if err == nil {
		err = h.DB.Files.UpdateDownloadCounter(c.Params("short"), file.Downloads)
		// check error, if error occours still send file if everything else works
		util.CheckWLogs(err)
		if file.Restricted {
			return c.Redirect("/r/" + c.Params("short"))
		}

		return c.Download("./data/uploads/" + file.Filename)
	}
	return c.Render("response", fiber.Map{
		"Text":        "Short does not exist",
		"Destination": "/",
	})
}

func (h handler) HandleDownloadRestricted(c *fiber.Ctx) error {
	file, err := h.DB.Files.GetFileByShort(c.Params("short"))
	if err == nil {
		return c.Download("./data/uploads/" + file.Filename)
	}
	return c.Render("response", fiber.Map{
		"Text":        "Short does not exist",
		"Destination": "/",
	})
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

		// ensure filenames are unique
		filenameUnique := false
		for !filenameUnique {
			_, err := os.Stat("./data/blind/" + file.Filename)

			if errors.Is(err, os.ErrNotExist) {
				filenameUnique = true
			} else if util.CheckWLogs(err) {
				return c.SendStatus(500)
			} else {
				file.Filename = "_" + file.Filename
			}
		}

		err = c.SaveFile(file, "./data/blind/"+file.Filename)

	case false:

		// ensure filenames are unique
		filenameUnique := false
		for !filenameUnique {
			_, err := h.DB.Files.GetFileByName(file.Filename)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				filenameUnique = true
			} else if util.CheckWLogs(err) {
				return c.SendStatus(500)
			} else {
				// ad an underscore to the filename if it's a duplicate
				file.Filename = "_" + file.Filename
			}
		}

		//edit
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
	short := c.Query("short", "")
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
	short := c.Query("short", "")
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

func (h handler) HandleSwitchRestrict(c *fiber.Ctx) error {
	short := c.Query("short", "")
	err := h.DB.Files.SwitchRestrict(short)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	return c.Redirect("/dashboard")
}

func (h handler) HandleDetails(c *fiber.Ctx) error {
	short := c.Query("short", "")
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
