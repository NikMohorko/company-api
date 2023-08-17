package main

import (
	"company_api/database"
	"company_api/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	database.Connect()
	app := fiber.New()

	setupPaths(app)

	err := app.Listen(":8080")

	if err != nil {
		log.Println(err.Error())
	}
}

func setupPaths(app *fiber.App) {

	// Users
	app.Post("/user/create", handlers.CreateUser)
	app.Post("/user/authenticate", handlers.Authenticate)

	// Companies
	app.Post("/company/create", handlers.CreateCompany)
	app.Patch("/company/update", handlers.UpdateCompany)
	app.Delete("/company/delete", handlers.DeleteCompany)
	app.Get("/company/get", handlers.GetCompany)

}
