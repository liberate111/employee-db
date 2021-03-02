package main

import (
	"emp/employees"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
)

func main() {
	engine := html.New("./templates", ".gohtml")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())
	app.Use(favicon.New())

	setupRoute(app)

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

func setupRoute(app *fiber.App) {
	app.Get("/", index)

	emp := app.Group("/emps")
	emp.Get("/", employees.Index)
	emp.Get("/show", employees.Show)
	emp.Get("/create", employees.Create)
	emp.Post("/create/process", employees.CreateProcess)
	emp.Get("/update", employees.Update)
	emp.Post("/update/process", employees.UpdateProcess)
	emp.Get("/delete/process", employees.DeleteProcess)
}

func index(c *fiber.Ctx) error {
	return c.Redirect("/emps")
}
