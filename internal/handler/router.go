package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/template/html"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/goupload/internal/database"
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
	engine := html.New("./web/templates", ".html")
	engine.AddFunc(
		"formatDate", func(timestamp time.Time) string {
			loc, _ := time.LoadLocation("Europe/Berlin")
			return timestamp.In(loc).Format("02.01.2006 15:04:05")
		},
	)
	app := fiber.New(fiber.Config{
		Views:       engine,
		BodyLimit:   uploadLimit * 1024 * 1024,
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	auth := jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/login?path=" + c.Path())
		},
		TokenLookup:   "cookie:JWT",
		SigningKey:    []byte(jwtkey),
		SigningMethod: "HS256",
	})

	app.Static("/static", "./web/public/static")
	app.Post("/login", handler.HandleLogin)
	app.Get("/login", template.Login)

	app.Get("/logout", handler.HandleLogout)

	app.Get("/", auth, template.Index)
	app.Post("/upload", auth, handler.HandleUploadWeb)
	app.Post("/", handler.HandleUploadSimple) // endpoint for simple upload, for example with curl
	app.Post("/paste", auth, handler.HandlePaste)

	dashboard := app.Group("/dashboard", auth)
	dashboard.Use(func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		if !claims["root"].(bool) {
			c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	})
	dashboard.Get("/", template.Dashboard)

	dashboard.Get("/removeFile/:short", handler.HandleRemoveFile)
	dashboard.Get("/moveToBlind/:short", handler.HandleMoveToBlind)
	dashboard.Get("/switchRestrict/:short", handler.HandleSwitchRestrict)
	dashboard.Get("/details/:short", handler.HandleDetails)

	dashboard.Get("/removePaste/:short", handler.HandleRemovePaste)

	dashboard.Get("/removeUser/:username", handler.HandleRemoveUser)

	dashboard.Get("/addUser", template.AddUser)
	dashboard.Post("/addUser", handler.HandleAddUser)

	dashboard.Get("/changePassword/:username", template.ChangePassword)
	dashboard.Post("/changePassword", handler.HandleChangePassword)

	dashboard.Get("/changePerms/:username", template.ChangePerms)
	dashboard.Post("/changePerms", handler.HandleChangePerms)

	dashboard.Get("/editPaste/:short", template.EditPaste)
	dashboard.Post("/editPaste", handler.HandleEditPaste)

	app.Get("/:short", handler.HandleShort)
	app.Get("/r/:short", auth, handler.HandleShort)

	app.Get("/raw/:short", handler.HandleShortsRaw)
	app.Get("/raw/r/:short", auth, handler.HandleShortsRaw)

	app.Listen(":" + port)
}
