package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/template/html"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/database"
)

type handler struct {
	DB     database.Database
	JwtKey string
	Url    string
}

func Start(port string, jwtkey string, url string, uploadLimit int, db database.Database) {
	handler := handler{db, jwtkey, url}
	template := template{
		DB:  db,
		Url: url,
	}
	engine := html.New("../data/templates", ".html")
	engine.AddFunc(
		"formatDate", func(timestamp time.Time) string {
			loc, _ := time.LoadLocation("Europe/Berlin")
			return timestamp.In(loc).Format("02.01.2006 15:04:05")
		},
	)
	app := fiber.New(fiber.Config{
		Views:                   engine,
		BodyLimit:               uploadLimit * 1024 * 1024,
		EnableTrustedProxyCheck: true,
	})

	auth := jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/login?path=" + c.Path())
		},
		TokenLookup:   "cookie:JWT",
		SigningKey:    []byte(jwtkey),
		SigningMethod: "HS256",
	})

	app.Static("/static", "../data/public/static")
	app.Post("/login", handler.HandleLogin)
	app.Get("/login", template.Login)

	app.Get("/logout", handler.HandleLogout)

	app.Get("/", auth, template.Index)
	app.Post("/upload", auth, handler.HandleUploadWeb)
	app.Post("/", handler.HandleUploadSimple) // endpoint for simple upload, for example with curl

	manage := app.Group("/manage", auth)
	manage.Use(func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		if claims["root"].(bool) != true {
			c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	})
	manage.Get("/", template.Manage)

	manage.Get("/removeFile", handler.HandleRemoveFile)
	manage.Get("/moveToBlind", handler.HandleMoveToBlind)
	manage.Get("/switchRestrict", handler.HandleSwitchRestrict)
	manage.Get("/details", handler.HandleDetails)
	manage.Post("/filterUser", template.HandleFilter) // unique template, handled by template instead of handler. TODO: change that maybe

	manage.Get("/removeUser", handler.HandleRemoveUser)
	manage.Post("/addUser", handler.HandleAddUser)
	manage.Post("/changePassword", handler.HandleChangePassword)

	app.Get("/:short", handler.HandleDownload)
	app.Get("/r/:short", auth, handler.HandleDownloadRestricted)

	app.Listen(":" + port)
}
