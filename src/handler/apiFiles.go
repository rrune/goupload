package handler

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/models"
	"github.com/rrune/goupload/util"
)

func (h handler) HandleDownload(c *fiber.Ctx) error {
	file, err := h.DB.Files.GetFileByShort(c.Params("short"))
	if err == nil {
		err = h.DB.Files.UpdateDownloadCounter(c.Params("short"), file.Downloads)
		// check error, if error occours still send file if everything else works
		util.CheckWLogs(err)
		if file.Restricted {
			return c.Redirect("/r/" + c.Params("short"))
		}

		return c.Download("../data/uploads/" + file.File)
	}
	return c.Render("response", fiber.Map{
		"Text": "Short does not exist",
	})
}

func (h handler) HandleDownloadRestricted(c *fiber.Ctx) error {
	file, err := h.DB.Files.GetFileByShort(c.Params("short"))
	if err == nil {
		return c.Download("../data/uploads/" + file.File)
	}
	return c.Render("response", fiber.Map{
		"Text": "Short does not exist",
	})
}

func (h handler) Upload(c *fiber.Ctx, username string, blindPerms bool, restrictedPerms bool, onetime bool, callback func(filename string, short string, blind bool) error) error {
	blind := c.FormValue("blind") == "blind"
	if blind && !blindPerms { // if blind is requested but not permitted
		return c.SendStatus(405)
	}

	restricted := c.FormValue("restricted") == "restricted"
	if restricted && !restrictedPerms { // if restricted is requested but not permitted
		return c.SendStatus(405)
	}

	file, err := c.FormFile("file")
	if util.CheckWLogs(err) {
		return c.SendString("No file given")
	}

	var short string
	switch blind {
	case true:
		err = c.SaveFile(file, "../data/blind/"+file.Filename)

	case false:
		short, err = h.DB.Files.AddNewFile(models.File{
			File:       file.Filename,
			Author:     username,
			Timestamp:  time.Now(),
			Ip:         c.IP(),
			Restricted: restricted,
		})
		if err == nil {
			err = c.SaveFile(file, "../data/uploads/"+file.Filename)
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
				"Text": "Uploaded " + filename,
			})
		}
		return c.Render("response", fiber.Map{
			"Text": "Uploaded " + filename,
			"Link": h.Url + short,
		})
	})
}

// simpler upload to use with curl
func (h handler) HandleUploadSimple(c *fiber.Ctx) error {
	// uses simple auth, there are no other permission checks prior to this one
	authBase64 := strings.Split(c.Get("Authorization"), " ")

	// The slice's length needs to be 2, otherwise the header was empty
	if len(authBase64) < 2 {
		c.SendString("No Authorization given")
	}
	// [0] of slice is "Basic"
	authByte, err := base64.StdEncoding.DecodeString(authBase64[1])
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	//basic auth is formated username:password
	auth := strings.Split(string(authByte), ":")
	if len(auth) != 2 {
		return c.SendString("Malformed Authorization header")
	}
	username := auth[0]
	password := auth[1]

	correct, user, err := h.ValidateCredentials(username, password)
	if util.CheckWLogs(err) {
		c.SendStatus(500)
	}
	if !correct {
		return c.SendString("Wrong Authorization")
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

		err = os.Remove("../data/uploads/" + file.File)
		if err == nil {

			err = h.DB.Files.RemoveFileByShort(short)
			if err == nil {
				return c.Render("response", fiber.Map{
					"Text": "Removed " + short,
				})
			}
		}
	}
	util.CheckWLogs(err)
	return c.Render("response", fiber.Map{
		"Text": "Could not remove " + short,
	})
}

func (h handler) HandleMoveToBlind(c *fiber.Ctx) error {
	short := c.Query("short", "")
	file, err := h.DB.Files.GetFileByShort(short)
	if err == nil {
		err = h.DB.Files.RemoveFileByShort(short)
		if err == nil {
			err = os.Rename("../data/uploads/"+file.File, "../data/blind/"+file.File)
		}
	}
	if util.CheckWLogs(err) {
		return c.Render("response", fiber.Map{
			"Text": "Could not move " + short,
		})
	}
	return c.Render("response", fiber.Map{
		"Text": "Moved " + short + "to blind",
	})
}

func (h handler) HandleSwitchRestrict(c *fiber.Ctx) error {
	short := c.Query("short", "")
	err := h.DB.Files.SwitchRestrict(short)
	if util.CheckWLogs(err) {
		return c.Render("response", fiber.Map{
			"Text": "Could restrict/unrestrict " + short,
		})
	}
	return c.Render("response", fiber.Map{
		"Text": "Restricted/Unrestriced " + short,
	})
}

func (h handler) HandleDetails(c *fiber.Ctx) error {
	short := c.Query("short", "")
	file, err := h.DB.Files.GetFileByShort(short)
	var info os.FileInfo
	if err == nil {
		info, err = os.Stat("../data/uploads/" + file.File)
	}
	if util.CheckWLogs(err) {
		return c.Render("response", fiber.Map{
			"Text": "Could not get details of " + short,
		})
	}
	infostrings := []string{
		info.Name(),
		strconv.FormatFloat(float64(info.Size())/1000, 'E', -1, 64) + " kb",
		strconv.Itoa(file.Downloads) + " Downloads",
	}
	return c.Render("details", fiber.Map{
		"Strings": infostrings,
	})
}
