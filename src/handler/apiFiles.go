package handler

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/models"
	"github.com/rrune/goupload/util"
)

func (h handler) Download(c *fiber.Ctx) error {
	file, err := h.DB.Files.GetFileByShort(c.Params("short"))
	if err == nil {
		if file.Restricted {
			return c.Redirect("/r/" + c.Params("short"))
		}

		return c.Download("../data/uploads/" + file.File)
	}
	return c.Render("response", fiber.Map{
		"Text": "Short does not exist",
	})
}

func (h handler) DownloadRestricted(c *fiber.Ctx) error {
	file, err := h.DB.Files.GetFileByShort(c.Params("short"))
	if err == nil {
		return c.Download("../data/uploads/" + file.File)
	}
	return c.Render("response", fiber.Map{
		"Text": "Short does not exist",
	})
}

func (h handler) Upload(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	var short string
	file, err := c.FormFile("file")

	restricted := false
	restrictedStr := c.FormValue("restricted")
	if restrictedStr == "restricted" {
		restricted = true
	}

	blind := c.FormValue("blind")
	if err == nil {
		if blind == "blind" {
			err = c.SaveFile(file, "../data/blind/"+file.Filename)
		} else {
			short, err = h.DB.Files.AddNewFile(models.File{
				File:       file.Filename,
				Author:     claims["username"].(string),
				Timestamp:  time.Now(),
				Ip:         c.IP(),
				Restricted: restricted,
			})
			if err == nil {
				err = c.SaveFile(file, "../data/uploads/"+file.Filename)
			}
		}
		if claims["onetime"].(bool) {
			err = h.DB.Users.RemoveUserByUsername(claims["username"].(string))
			c.ClearCookie("JWT")
		}
	}
	if util.Check(err) {
		return c.Render("response", fiber.Map{
			"Text": "Could not upload your file",
		})
	}
	if blind == "blind" {
		return c.Render("response", fiber.Map{
			"Text": "Uploaded " + file.Filename,
		})
	}
	return c.Render("response", fiber.Map{
		"Text": "Uploaded " + file.Filename,
		"Link": h.Url + short,
	})
}

func (h handler) RemoveFile(c *fiber.Ctx) error {
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
	return c.Render("response", fiber.Map{
		"Text": "Could not remove " + short,
	})
}

func (h handler) MoveToBlind(c *fiber.Ctx) error {
	short := c.Query("short", "")
	file, err := h.DB.Files.GetFileByShort(short)
	if err == nil {
		err = h.DB.Files.RemoveFileByShort(short)
		if err == nil {
			err = os.Rename("../data/uploads/"+file.File, "../data/blind/"+file.File)
		}
	}
	if util.Check(err) {
		return c.Render("response", fiber.Map{
			"Text": "Could not move " + short,
		})
	}
	return c.Render("response", fiber.Map{
		"Text": "Moved " + short + "to blind",
	})
}

func (h handler) Details(c *fiber.Ctx) error {
	short := c.Query("short", "")
	file, err := h.DB.Files.GetFileByShort(short)
	var info os.FileInfo
	if err == nil {
		info, err = os.Stat("../data/uploads/" + file.File)
	}
	if util.Check(err) {
		return c.Render("response", fiber.Map{
			"Text": "Could not get details of " + short,
		})
	}
	infostr := info.Name() + " " + strconv.FormatInt(info.Size(), 10) + " bytes"
	return c.Render("response", fiber.Map{
		"Text": infostr,
	})
}
