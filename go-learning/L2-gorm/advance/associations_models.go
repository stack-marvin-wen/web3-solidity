package advance

import "time"

type user struct {
	Id        uint
	Name      string
	Email     string
	Profile   profile // 一对一关系
	Orders    []order
	Roles     []role `gorm:"many2many:user_roles;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type profile struct {
	Id        int
	UserID    uint `gorm:"uniqueIndex"` // 外键，唯一索引
	Nickname  string
	Phone     string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type order struct {
	Id         int
	UserID     uint
	OrderNo    string `gorm:"uniqueIndex"`
	Items      []orderitem
	TotalPrice float64
	status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type orderitem struct {
	Id        int
	OrderID   uint
	ProductID uint
	Product   product
	UnitPrice int64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
type product struct {
	Id        int
	Name      string
	Price     int64
	SKU       string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type role struct {
	Id          int
	Name        string `gorm:"uniqueIndex"`
	Description string
	Users       []user `gorm:"many2many:user_roles;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
