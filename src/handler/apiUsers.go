package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/models"
	"github.com/rrune/goupload/util"
)

func (h handler) HandleLogin(c *fiber.Ctx) error {
	CallbackPath := c.Query("path", "")

	l := new(models.Login)
	if err := c.BodyParser(l); err != nil {
		return err
	}
	user, err := h.DB.Users.GetUserByUsername(l.Username)
	if util.Check(err) {
		return c.Redirect("/login?msg=Wrong username or password&path=" + CallbackPath)
	}
	correct := util.DoPasswordsMatch(user.Password, l.Password)

	if correct {
		claims := jwt.MapClaims{
			"username":   user.Username,
			"root":       user.Root,
			"blind":      user.Blind,
			"onetime":    user.Onetime,
			"restricted": user.Restricted,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(h.JwtKey))
		if util.Check(err) {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		cookie := fiber.Cookie{
			Name:  "JWT",
			Value: t,

			HTTPOnly: true,
		}
		c.Cookie(&cookie)

		return c.Redirect(CallbackPath)
	} else {
		return c.Redirect("/login?msg=Wrong username or password&path=" + CallbackPath)
	}
}

func (h handler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("JWT")
	return c.Redirect("/login")
}

func (h handler) AddUser(c *fiber.Ctx) error {
	formUser := new(models.UserFromForm)
	if err := c.BodyParser(formUser); err != nil {
		return err
	}

	user := models.User{
		Username:   formUser.Username,
		Password:   formUser.Password,
		Root:       formUser.Root == "root",
		Blind:      formUser.Blind == "blind",
		Onetime:    formUser.Onetime == "onetime",
		Restricted: formUser.Onetime == "restricted",
	}

	err := h.DB.Users.CreateUser(&user)
	if util.Check(err) {
		return c.Render("response", fiber.Map{
			"Text": "Could not create " + user.Username,
		})
	}
	return c.Render("response", fiber.Map{
		"Text": "Created " + user.Username,
	})
}

func (h handler) RemoveUser(c *fiber.Ctx) error {
	username := c.Query("username", "")
	err := h.DB.Users.RemoveUserByUsername(username)
	if util.Check(err) {
		return c.Render("response", fiber.Map{
			"Text": "Could not remove " + username,
		})
	}
	return c.Render("response", fiber.Map{
		"Text": "Removed " + username,
	})
}

func (h handler) ChangePassword(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	err := h.DB.Users.ChangePassword(user.Username, user.Password)
	if util.Check(err) {
		return c.Render("response", fiber.Map{
			"Text": "Could not change password of " + user.Username,
		})
	}
	return c.Render("response", fiber.Map{
		"Text": "Changed Password of " + user.Username,
	})
}
