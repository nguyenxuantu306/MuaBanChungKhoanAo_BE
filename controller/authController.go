package controller

import (
	"mymodule/database"
	"mymodule/models"

	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const Secret = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Kiểm tra xem email đã tồn tại trong cơ sở dữ liệu hay chưa
	existingUser := models.Users{}
	database.DB.Where("email = ?", data["email"]).First(&existingUser)

	if existingUser.UserId != 0 {
		// Nếu email đã tồn tại, trả về thông báo lỗi
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), int(12))
	users := models.Users{
		UserName: data["name"],
		Email:    data["email"],
		Password: string(password),
	}

	database.DB.Create(&users)

	return c.JSON(users)
}

func generateAccessToken(user models.Users) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.UserId)),
		ExpiresAt: time.Now().Add(time.Hour).Unix(), // Set your preferred expiration time
	})

	token, err := claims.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateRefreshToken(user models.Users) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.UserId)),
		ExpiresAt: time.Now().Add(5 * time.Second).Unix(), // Set your preferred expiration time
	})

	token, err := claims.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.FormValue("refresh_token")

	// Xác minh refresh token
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claims := token.Claims.(jwt.MapClaims)

	// Tạo mới access token
	userID, _ := strconv.Atoi(claims["iss"].(string))
	user := models.Users{UserId: uint(userID)}

	accessToken, accessErr := generateAccessToken(user)
	if accessErr != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Could not issue a new access token",
		})
	}

	// Đặt access token mới vào cookie
	accessCookie := fiber.Cookie{
		Name:     "jwt",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour), // Đặt thời gian hết hạn accessToken
		HTTPOnly: true,
	}
	c.Cookie(&accessCookie)

	return c.JSON(fiber.Map{
		"message": "Access token refreshed",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.Users

	// Tìm người dùng dựa trên địa chỉ email
	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.UserId == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// So sánh mật khẩu được cung cấp với mật khẩu trong cơ sở dữ liệu
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	// Tạo mã access token mới
	accessToken, accessErr := generateAccessToken(user)
	if accessErr != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Could not issue an access token",
		})
	}

	// Tạo mã refresh token mới
	refreshToken, refreshErr := generateRefreshToken(user)
	if refreshErr != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Could not issue a refresh token",
		})
	}

	// Đặt mã access token vào cookie
	accessCookie := fiber.Cookie{
		Name:     "jwt",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&accessCookie)

	return c.JSON(fiber.Map{
		"message":       "Success",
		"refresh_token": refreshToken,
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.Users

	// Tìm người dùng dựa trên ID trong claims
	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Second),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}
