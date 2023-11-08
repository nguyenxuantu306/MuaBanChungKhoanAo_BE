package models

import "time"

// User model
type Users struct {
	UserId         uint    `gorm:"not null;primary_key" json:"userId" autoIncrement`
	UserName       string  `gorm:"not null" json:"userName"`
	Password       string  `gorm:"not null" json:"password"`
	FullName       string  `json:"fullName"`
	Email          string  `json:"email"`
	PhoneNumber    string  `json:"phoneNumber"`
	AccountBalance float64 `gorm:"type:decimal(10,3);" json:"accountBalance"`
}

// CommandType model
type CommandTypes struct {
	CommandTypeId   string `gorm:"not null;primary_key" json:"commandTypeId"`
	CommandTypeName string `gorm:"not null" json:"commandTypeName"`
}

// Stock model
type Stocks struct {
	StockId    string    `gorm:"not null;primary_key" json:"stockId"`
	StockName  string    `gorm:"not null" json:"stockName"`
	OpenTime   time.Time `json:"openTime"`
	OpenPrice  float64   `gorm:"type:decimal(10,3);" json:"openPrice"`
	CloseTime  time.Time `json:"closeTime"`
	ClosePrice float64   `gorm:"type:decimal(10,3);" json:"closePrice"`
	Price      float64   `gorm:"type:decimal(10,3);" json:"price"`
}

// UserInfo model
type UserInfos struct {
	UserInfoId    uint    `gorm:"not null;primary_key" json:"userInfoId" autoIncrement`
	UserId        uint    `gorm:"not null" json:"userId"`
	StockId       string  `gorm:"not null" json:"stockId"`
	StockName     string  `gorm:"not null" json:"stockName"`
	StockQuantity float64 `json:"stockQuantity"`
}

// Order model
type Orders struct {
	OrderId            uint      `gorm:"not null;primary_key" json:"orderId" autoIncrement`
	Price              float64   `gorm:"type:decimal(10,3);" json:"price"`
	Quantity           float64   `json:"quantity"`
	Total              float64   `json:"total"`
	ImplementationDate time.Time `json:"implementationDate"`
	StockId            string    `gorm:"not null" json:"stockId"`
	CommandTypeId      string    `gorm:"not null" json:"commandTypeId"`
	UserId             uint      `gorm:"not null" json:"userId"`
	Command            string    `gorm:"not null" json:"command"`
	Status             string    `gorm:"not null" json:"status"`
}
