package handlers

import (
	"company_api/database"
	"company_api/models"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func validateCompany(company *models.Company) error {

	if company.Name == "" {
		return errors.New("Company name is required.")
	}

	if len(company.Name) > 15 {
		return errors.New("Company name must be between 1 and 15 characters.")
	}

	if len(company.Description) > 3000 {
		return errors.New("Description must be shorter than 3000 characters.")
	}

	if company.EmployeeCount == nil {
		return errors.New("Field \"employee_count\" is required.")
	}

	if company.Registered == nil {
		return errors.New("Field \"registered\" is required.")
	}
	if company.Type == "" {
		return errors.New("Field \"type\" is required.")
	}

	allowedTypes := []string{"Corporations", "NonProfit", "Cooperative", "Sole Proprietorship"}

	if !slices.Contains(allowedTypes, company.Type) {
		return errors.New("Invalid company type. Allowed values are: " + strings.Join(allowedTypes, ", ") + ".")
	}

	return nil
}

func CreateCompany(c *fiber.Ctx) error {

	tokenValid, errMessage := validateToken(c)

	if !tokenValid {
		return c.Status(fiber.StatusUnauthorized).JSON(errMessage)
	}

	company := new(models.Company)

	// Automatically create UUID for new companies
	company.ID = uuid.NewString()

	if err := c.BodyParser(company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body contains invalid data types.",
		})
	}

	// Validate company data
	if err := validateCompany(company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if creation := database.DB.Db.Create(&company); creation.Error != nil {

		// Postgres unique value violation
		if errors.Is(creation.Error, gorm.ErrDuplicatedKey) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Company name already exists.",
			})

		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": creation.Error.Error()})

		}
	}

	return c.Status(201).JSON(company)
}

func GetCompany(c *fiber.Ctx) error {

	// Read URL params
	params := c.Queries()

	companyName, ok := params["name"]

	company := new(models.Company)

	if ok {

		if query := database.DB.Db.First(&company, "name = ?", companyName); query.Error != nil {

			if errors.Is(query.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Company name does not exist.",
				})

			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(nil)
			}
		}

		return c.Status(200).JSON(company)

	} else {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing company name in URL parameters.",
		})
	}

}

func UpdateCompany(c *fiber.Ctx) error {

	tokenValid, errMessage := validateToken(c)

	if !tokenValid {
		return c.Status(fiber.StatusUnauthorized).JSON(errMessage)
	}

	company := new(models.Company)

	if err := c.BodyParser(company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if update := database.DB.Db.Model(&company).Where("name = ?", company.Name).Updates(company); update.Error != nil {

		if errors.Is(update.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Company name does not exist.",
			})

		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(nil)
		}

	}

	return c.Status(200).JSON(company)

}

func DeleteCompany(c *fiber.Ctx) error {

	tokenValid, errMessage := validateToken(c)

	if !tokenValid {
		return c.Status(fiber.StatusUnauthorized).JSON(errMessage)
	}

	params := c.Queries()

	companyName, ok := params["name"]

	company := new(models.Company)

	if ok {
		// Use unscoped to prevent soft deletion - violates unique contraint if you want to create the same company again
		if deletion := database.DB.Db.Unscoped().Where("name = ?", companyName).Delete(&company); deletion.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(nil)
		}

		return c.Status(200).JSON(nil)

	} else {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Company Name not provided.",
		})
	}

}
