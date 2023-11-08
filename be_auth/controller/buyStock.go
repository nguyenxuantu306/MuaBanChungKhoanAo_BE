package controller

import (
	"errors"
	"fmt"
	"mymodule/database"
	"mymodule/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BuyStock(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		c.Status(fiber.StatusAccepted)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)
	user_idStr := claims.Issuer
	user_idInt, err := strconv.Atoi(user_idStr)
	user_id := uint(user_idInt)
	if err != nil {
		return err
	}

	stock_id := data["stockid"]
	quantityStr, quantityExists := data["quantity"]
	priceStr, priceExists := data["price"]

	var quantity float64
	var price float64
	if quantityExists {
		quantity, err = strconv.ParseFloat(quantityStr, 64)
		if err != nil {
			return err
		}
	}
	if priceExists {
		price, err = strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return err
		}
	}

	user := models.Users{}
	userInfo := models.UserInfos{}

	stock := models.Stocks{}
	if err := database.DB.Where("stock_id = ?", stock_id).First(&stock).Error; err != nil {
		return err
	}

	if err := database.DB.Where("user_id = ?", user_id).First(&user).Error; err != nil {
		return err
	}

	if err := database.DB.Where("user_id = ? AND stock_id = ?", user_id, stock_id).First(&userInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUserInfo := models.UserInfos{
				UserId:        user_id,
				StockId:       stock_id,
				StockName:     stock.StockName,
				StockQuantity: 0,
			}
			if err := database.DB.Create(&newUserInfo).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	stockPrice := stock.Price

	if priceExists && quantityExists {
		// Cả giá và số lượng đều được cung cấp, hãy tính tổng chi phí
		totalCost := price
		if user.AccountBalance < totalCost {
			return fmt.Errorf("Số tiền không đủ để mua cổ phiếu")
		}

		quantity = price / stockPrice // Tính số lượng dựa trên giá
	} else if priceExists {
		// Chỉ có giá được cung cấp, tính toán số lượng dựa trên giá
		quantity = price / stockPrice
	} else if quantityExists {
		// Chỉ cung cấp số lượng, tính tổng chi phí
		totalCost := float64(quantity) * stockPrice
		if user.AccountBalance < totalCost {
			return fmt.Errorf("Số tiền không đủ để mua cổ phiếu")
		}
		price = quantity * stockPrice
	}

	newAccountBalance := user.AccountBalance - price
	newStockQuantity := userInfo.StockQuantity + quantity

	if err := database.DB.Model(&user).Updates(models.Users{
		AccountBalance: newAccountBalance,
	}).Error; err != nil {
		return err
	}

	if err := database.DB.Model(&userInfo).Updates(models.UserInfos{
		StockQuantity: newStockQuantity,
	}).Error; err != nil {
		return err
	}

	command_type_id := data["commandtypeid"]
	command := data["command"]

	order := models.Orders{
		UserId:             user_id,
		StockId:            stock_id,
		CommandTypeId:      command_type_id,
		Command:            command,
		ImplementationDate: time.Now(),
		Quantity:           quantity,
		Price:              price,
		Total:              quantity * price,
	}

	result := database.DB.Create(&order)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(order)
}
