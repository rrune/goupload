package handler

import (
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
	app := fiber.New(fiber.Config{
		Views:     engine,
		BodyLimit: uploadLimit * 1024 * 1024,
	})

	auth := jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/login")
		},
		TokenLookup:   "cookie:JWT",
		SigningKey:    []byte(jwtkey),
		SigningMethod: "HS256",
	})

	app.Static("/static", "../data/public/static")
	app.Post("/login", handler.HandleLogin)
	app.Get("/login", template.Login)

	app.Get("/logout", handler.Logout)

	app.Get("/", auth, template.Index)
	app.Post("/upload", auth, handler.Upload)

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

	manage.Get("/removeFile", handler.RemoveFile)
	manage.Get("/moveToBlind", handler.MoveToBlind)
	manage.Get("/details", handler.Details)
	manage.Post("/filterUser", template.Filter)

	manage.Get("/removeUser", handler.RemoveUser)
	manage.Post("/addUser", handler.AddUser)
	manage.Post("/changePassword", handler.ChangePassword)

	app.Get("/:short", handler.Download)

	app.Listen(":" + port)
}
