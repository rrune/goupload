package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rrune/goupload/internal/util"
)

func (h handler) HandleShort(c *fiber.Ctx) error {
	short, err := h.DB.Shorts.GetShort(c.Params("short"))
	if err == nil {

		if c.OriginalURL()[:3] != "/r/" {
			err = h.DB.Shorts.UpdateDownloadCounter(c.Params("short"), short.Downloads)
			// check error, if error occours still send if everything else works
			util.CheckWLogs(err)

			if short.Restricted {
				return c.Redirect("/r/" + c.Params("short"))
			}
		}

		switch short.Type {
		case "file":
			return h.Download(c)
		case "paste":
			return h.Paste(c)
		}
	}
	return c.Render("response", fiber.Map{
		"Text":        "Short does not exist",
		"Destination": "/",
	})
}

func (h handler) HandleShortsRaw(c *fiber.Ctx) error {
	short, err := h.DB.Shorts.GetShort(c.Params("short"))
	if err == nil {

		if c.OriginalURL()[4:7] != "/r/" {
			err = h.DB.Shorts.UpdateDownloadCounter(c.Params("short"), short.Downloads)
			// check error, if error occours still send if everything else works
			util.CheckWLogs(err)

			if short.Restricted {
				return c.Redirect("/raw/r/" + c.Params("short"))
			}
		}

		switch short.Type {
		case "file":
			return c.Redirect("/" + c.Params("short"))
		case "paste":
			return h.GetPasteRaw(c)
		}
	}
	return c.Render("response", fiber.Map{
		"Text":        "Short does not exist",
		"Destination": "/",
	})
}

func (h handler) HandleSwitchRestrict(c *fiber.Ctx) error {
	short := c.Params("short", "")
	if short == "" {
		return c.SendStatus(400)
	}
	err := h.DB.Shorts.SwitchRestrict(short)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	return c.Redirect("/dashboard")
}
