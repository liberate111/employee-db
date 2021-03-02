package employees

import (
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.Status(405).SendString("StatusMethodNotAllowed")
	}

	emps, err := AllEmps()
	if err != nil {
		return c.Status(500).SendString("StatusInternalServerError")
	}
	return c.Render("employees", emps)
}

func Show(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.Status(405).SendString("StatusMethodNotAllowed")
	}
	emp, err := OneEmp(c)
	if err != nil {
		return c.Status(500).SendString("StatusInternalServerError")
	}
	return c.Render("show", emp)
}

func Create(c *fiber.Ctx) error {
	return c.Render("create", nil)
}

func CreateProcess(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return c.Status(405).SendString("StatusMethodNotAllowed")
	}

	emp, err := PutEmp(c)
	if err != nil {
		return c.Status(500).SendString("StatusInternalServerError")
	}

	return c.Render("created", emp)
}

func Update(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.Status(405).SendString("StatusMethodNotAllowed")
	}

	emp, err := OneEmp(c)
	if err != nil {
		return c.Status(500).SendString("StatusInternalServerError")
	}

	return c.Render("update", emp)
}

func UpdateProcess(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return c.Status(405).SendString("StatusMethodNotAllowed")
	}

	emp, err := UpdateEmp(c)
	if err != nil {
		return c.Status(500).SendString("StatusInternalServerError")
	}

	return c.Render("updated", emp)
}

func DeleteProcess(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.Status(405).SendString("StatusMethodNotAllowed")
	}

	err := DeleteEmp(c)
	if err != nil {
		return c.Status(400).SendString("StatusBadRequest")
	}
	return c.Redirect("/emps")
}
