package advance

import (
	"L2-gorm/testutil"
	"testing"

	"gorm.io/gorm"
)

func TestTransaction(t *testing.T) {
	db := testutil.NewTestDB(t, "transaction.db")
	db.AutoMigrate(&user{}, &profile{}, &order{}, &orderitem{}, &product{}, &role{})

	users := []user{
		{Name: "user1", Email: "user1@example.com", Profile: profile{
			Nickname: "user1",
			Phone:    "1234567890121",
			Address:  "Beijing, China",
		},
			Orders: []order{
				{OrderNo: "order1", TotalPrice: 100.0, status: "pending", Items: []orderitem{
					{Product: product{Name: "product1", Price: 50.0, SKU: "sku1"}, UnitPrice: 50.0, Quantity: 2},
				}},
			},
			Roles: []role{
				{Name: "admin", Description: "管理员"},
				{Name: "user", Description: "普通用户"},
			},
		},
		{Name: "user2", Email: "user2@example.com", Profile: profile{
			Nickname: "user2",
			Phone:    "1234567890122",
			Address:  "Shanghai, China",
		},
			Orders: []order{
				{
					OrderNo: "order2", TotalPrice: 200.0, status: "pending", Items: []orderitem{
						{Product: product{Name: "product2", Price: 100.0, SKU: "sku2"}, UnitPrice: 100.0, Quantity: 2},
					},
				},
			},
			Roles: []role{
				{Name: "user", Description: "普通用户"},
			},
		},
	}
	db.Session(&gorm.Session{}).Create(&users)
	t.Run("transaction/auto", func(t *testing.T) {
		err := db.Transaction(func(tx *gorm.DB) error {
			user := user{Name: "user3", Email: "user3@example.com", Profile: profile{
				Nickname: "user3",
				Phone:    "1234567890123",
				Address:  "Guangzhou, China",
			},
				Orders: []order{
					{OrderNo: "order3", TotalPrice: 300.0, status: "pending", Items: []orderitem{
						{Product: product{Name: "product3", Price: 150.0, SKU: "sku"}, UnitPrice: 150.0, Quantity: 2},
					}},
				},
				Roles: []role{
					{Name: "user", Description: "普通用户"},
				},
			}
			err := tx.Session(&gorm.Session{}).Create(&user).Error
			if err != nil {
				return err
			}
			return nil

		})
		if err != nil {
			t.Fatal(err)
		}
	})

}
