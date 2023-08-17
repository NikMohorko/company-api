package handlers

import (
	"company_api/database"
	"company_api/models"
	"errors"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTData struct {
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom_claims"`
}

func generateToken(username string) (string, int64, error) {

	validity := time.Now().Add(20 * time.Minute).Unix()

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: validity,
		},
		CustomClaims: map[string]string{
			"username": username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)

	return tokenString, validity, err

}

func validateToken(c *fiber.Ctx) (isValid bool, message interface{}) {

	requestHeaders := c.GetReqHeaders()
	claims := JWTData{}

	requestToken, ok := requestHeaders["Authorization"]

	if !ok {
		return false, nil
	}

	splitToken := strings.Split(requestToken, "Bearer ")

	if len(splitToken) != 2 {
		return false, nil

	} else {
		requestToken = splitToken[1]
	}

	token, err := jwt.ParseWithClaims(requestToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {

		if !token.Valid {
			return false, "Token expired."

		} else {
			return false, nil
		}

	} else {
		return true, nil
	}

}

func validateUser(user *models.User) error {

	if user.Username == "" {
		return errors.New("Username is required.")
	}

	if len(user.Username) < 3 || len(user.Username) > 50 {
		return errors.New("Username must be between 3 and 50 characters long.")
	}

	if user.Password == "" {
		return errors.New("Password is required.")
	}

	if len(user.Password) < 6 || len(user.Password) > 30 {
		return errors.New("Password must be between 6 and 30 characters long.")
	}

	return nil

}

func CreateUser(c *fiber.Ctx) error {

	user := new(models.User)

	user.ID = uuid.NewString()

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body contains invalid data types.",
		})
	}

	// Validate user data
	if err := validateUser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 6)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Store hashed password in db
	user.Password = string(hashedPassword)

	if creation := database.DB.Db.Create(&user); creation.Error != nil {

		if errors.Is(creation.Error, gorm.ErrDuplicatedKey) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Username already exists.",
			})

		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": creation.Error.Error()})

		}
	}

	return c.Status(201).JSON(nil)

}

func Authenticate(c *fiber.Ctx) error {

	receivedData := new(models.User)

	if err := c.BodyParser(receivedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	storedData := new(models.User)

	// Check if user exists
	if getUser := database.DB.Db.First(&storedData, "username = ?", receivedData.Username); getUser.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(nil)
	}

	// Check if stored hashed password matches the provided one
	err := bcrypt.CompareHashAndPassword([]byte(storedData.Password), []byte(receivedData.Password))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(nil)
	}

	token, validity, err := generateToken(receivedData.Username)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nil)
	}

	token_response := map[string]interface{}{
		"token":       token,
		"valid_until": validity,
	}

	return c.Status(200).JSON(token_response)

}
