package controller

import (
	"fmt"
	"mymodule/database"
	"mymodule/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func SellStock(c *fiber.Ctx) error {
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
	command_type_id := data["commandtypeid"]
	command := data["command"]

	stock := models.Stocks{}
	if err := database.DB.Where("stock_id = ?", stock_id).First(&stock).Error; err != nil {
		return err
	}

	quantityStr, quantityExists := data["quantity"]
	purchasePriceStr, purchasePriceExists := data["price"]

	if !quantityExists && !purchasePriceExists {
		// Không có dữ liệu nào được cung cấp, xử lý lỗi tùy ý
		return fmt.Errorf("Dữ liệu thiếu")
	}

	var quantity float64
	if quantityExists {
		quantity, err = strconv.ParseFloat(quantityStr, 64)
		if err != nil {
			return err
		}
	}

	var price float64
	if purchasePriceExists {
		price, err = strconv.ParseFloat(purchasePriceStr, 64)
		if err != nil {
			return err
		}
	} else {
		price = stock.Price
	}

	var total float64
	if quantityExists && purchasePriceExists {
		total = quantity * price
	} else if quantityExists {
		total = quantity * stock.Price
	} else {
		fmt.Errorf("Số lượng hoặc giá mua cổ phiếu không được cung cấp")
	}

	user := models.Users{}
	if err := database.DB.Where("user_id = ?", user_id).First(&user).Error; err != nil {
		return err
	}

	userInfo := models.UserInfos{}
	// 2. Kiểm tra xem người dùng có đủ số lượng cổ phiếu để bán không
	if err := database.DB.Where("user_id = ? AND stock_id = ?", user_id, stock_id).First(&userInfo).Error; err != nil {
		return fmt.Errorf("Người dùng không sở hữu cổ phiếu này")
	}

	if userInfo.StockQuantity < quantity {
		return fmt.Errorf("Số lượng cổ phiếu không đủ để bán")
	}

	newAccountBalance := user.AccountBalance + (float64(quantity) * price)
	newStockQuantity := userInfo.StockQuantity - quantity

	if err := database.DB.Model(&userInfo).Updates(models.UserInfos{
		StockQuantity: newStockQuantity,
	}).Error; err != nil {
		return err
	}

	if err := database.DB.Model(&user).Updates(models.Users{
		AccountBalance: newAccountBalance,
	}).Error; err != nil {
		return err
	}

	if userInfo.StockQuantity == 0 {
		if err := database.DB.Delete(&userInfo).Error; err != nil {
			return err
		}
	}

	order := models.Orders{
		UserId:             user_id,
		StockId:            stock_id,
		CommandTypeId:      command_type_id,
		Command:            command,
		ImplementationDate: time.Now(),
		Quantity:           quantity,
		Price:              price,
		Total:              total,
	}

	result := database.DB.Create(&order)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(order)
}
