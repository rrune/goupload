package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/internal/models"
	"github.com/rrune/goupload/internal/util"
)

func (h handler) ValidateCredentials(username string, password string) (bool, models.User, error) {
	user, err := h.DB.Users.GetUserByUsername(username)
	if util.Check(err) {
		return false, models.User{}, err
	}

	correct := util.DoPasswordsMatch(user.Password, password)
	if correct {
		return true, user, err
	}
	return false, models.User{}, err
}

func (h handler) HandleLogin(c *fiber.Ctx) error {
	CallbackPath := c.Query("path", "")
	if CallbackPath == "" {
		CallbackPath = "/"
	}

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
		if util.CheckWLogs(err) {
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

func (h handler) HandleLogout(c *fiber.Ctx) error {
	c.ClearCookie("JWT")
	return c.Redirect("/login")
}

func (h handler) HandleAddUser(c *fiber.Ctx) error {
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
		Restricted: formUser.Restricted == "restricted",
	}

	err := h.DB.Users.CreateUser(&user)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	return c.Redirect("/dashboard")
}

func (h handler) HandleRemoveUser(c *fiber.Ctx) error {
	username := c.Params("username", "")
	if username == "" {
		return c.SendStatus(400)
	}
	err := h.DB.Users.RemoveUserByUsername(username)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	return c.Redirect("/dashboard")
}

func (h handler) HandleChangePassword(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	err := h.DB.Users.ChangePassword(user.Username, user.Password)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	return c.Render("response", fiber.Map{
		"Text":        "Changed Password of " + user.Username,
		"Destination": "/dashboard",
	})
}

func (h handler) HandleChangePerms(c *fiber.Ctx) error {
	formUser := new(models.UserFromForm)
	if err := c.BodyParser(formUser); err != nil {
		return err
	}

	user := models.User{
		Username:   formUser.Username,
		Password:   "",
		Root:       formUser.Root == "root",
		Blind:      formUser.Blind == "blind",
		Onetime:    formUser.Onetime == "onetime",
		Restricted: formUser.Restricted == "restricted",
	}
	err := h.DB.Users.ChangePerms(user)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	return c.Redirect("/dashboard")
}
