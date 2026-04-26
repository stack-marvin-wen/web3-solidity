package basic

import (
	"L2-gorm/testutil"
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestCRUD(t *testing.T) {
	db := testutil.NewTestDB(t, "crud.db")
	type User struct {
		ID        uint      `gorm:"column:id;type:int;primaryKey;autoIncrement"`
		Name      string    `gorm:"size:64;not null"`
		Email     string    `gorm:"size:128;unique;not null"`
		Age       int       `gorm:"type:int;not null"`
		Status    string    `gorm:"size:16;default:active;index"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("auto migrate err: %v", err)
	}
	if err := db.Exec("delete from users").Error; err != nil {
		t.Fatalf("delete err: %v", err)
	}

	users := []User{
		{Name: "Alice", Email: "alice@example.com", Age: 18, Status: "active"},
		{Name: "Bob", Email: "bob@example.com", Age: 19, Status: "active"},
		{Name: "Charlie", Email: "charlie@example.com", Age: 20, Status: "inactive"},
		{Name: "David", Email: "david@example.com", Age: 21, Status: "active"},
		{Name: "Edward", Email: "edward@example.com", Age: 22, Status: "active"},
		{Name: "Frank", Email: "frank@example.com", Age: 23, Status: "inactive"},
		{Name: "George", Email: "george@example.com", Age: 24, Status: "active"},
		{Name: "Helen", Email: "helen@example.com", Age: 25, Status: "active"},
		{Name: "Irene", Email: "irene@example.com", Age: 26, Status: "active"},
		{Name: "James", Email: "james@example.com", Age: 27, Status: "active"},
		{Name: "Kevin", Email: "kevin@example.com", Age: 28, Status: "active"},
		{Name: "Lucas", Email: "lucas@example.com", Age: 29, Status: "active"},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatalf("create user err: %v", err)
	}

	// single user create test
	t.Run("create", func(t *testing.T) {
		u := User{Name: "Mike", Email: "mike@example.com", Age: 30, Status: "active"}
		if err := db.Create(&u).Error; err != nil {
			t.Fatalf("create user err: %v", err)
		}
		t.Logf("new user id=%d", u.ID)
	})
	// query test
	t.Run("query/first", func(t *testing.T) {
		var user User
		if err := db.Where("status = ?", "active").First(&user).Error; err != nil {
			t.Fatalf("query first active user err: %v", err)
		}
		t.Logf("first active user: %+v", user)
	})
	t.Run("query/take", func(t *testing.T) {
		// will take any random user, no order by default
		var user User
		if err := db.Take(&user).Error; err != nil {
			t.Fatalf("query take user err: %v", err)
		}
		t.Logf("take user: %+v", user)
	})
	t.Run("query/find", func(t *testing.T) {
		var users []User
		if err := db.Where("age >= ?", 25).Order("created_at desc").Find(&users).Error; err != nil {
			t.Fatalf("query find users err: %v", err)
		}
		t.Logf("users with age >= 25: %+v", users)
		var allUsers []User
		if err := db.Find(&allUsers).Error; err != nil {
			t.Fatalf("query find users err: %v", err)
		}
		t.Logf("all users: %+v", allUsers)
	})
	t.Run("query/scan", func(t *testing.T) {
		type UserSummary struct {
			Name   string
			Email  string
			Status string
		}
		var summaries []UserSummary
		if err := db.Model(&User{}).Select("name, email, status").Where("age >= ?", 25).Scan(&summaries).Error; err != nil {
			t.Fatalf("query scan user summaries err: %v", err)
		}
		t.Logf("user summaries with age >= 25: %+v", summaries)
	})
	// update test
	t.Run("update", func(t *testing.T) {
		var user User
		if err := db.Where("email = ?", "alice@example.com").First(&user).Error; err != nil {
			t.Fatalf("query first active user err: %v", err)
		}
		fmt.Printf("old record: %v\n", user)
		if err := db.Model(&user).Select("Age", "Status").Where("email = ?", "alice@example.com").Updates(User{Age: 30, Status: "vip"}).Error; err != nil {
			t.Fatalf("update user err: %v", err)
		}
		if err := db.Where("email = ?", "alice@example.com").First(&user).Error; err != nil {
			t.Fatalf("query first active user err: %v", err)
		}
		fmt.Printf("new record(Age: 30, Status: vip): %v\n", user)

	})
	t.Run("bulk update", func(t *testing.T) {
		// Model(&User{}): Specify the model for bulk operation
		// Where: Add conditions to filter which records to update
		// Updates: Update all matching records
		// Using map[string]any allows updating specific fields without zero value issues
		res := db.Model(&User{}).Where("status = ?", "inactive").Updates(map[string]any{"status": "pending_review"})
		if res.Error != nil {
			t.Fatalf("bulk update: %v", res.Error)
		}
		// RowsAffected: Check how many rows were actually updated
		if res.RowsAffected == 0 {
			t.Fatalf("expected rows to be updated")
		}
	})
	// delete test
	t.Run("delete", func(t *testing.T) {
		var user User
		// First: Load the user to delete
		if err := db.First(&user, "email = ?", "alice1@example.com").Error; err != nil {
			t.Fatalf("load user: %v", err)
		}
		// Delete: Delete by primary key
		// First parameter is the model type, second is the primary key value
		if err := db.Delete(&User{}, user.ID).Error; err != nil {
			t.Fatalf("delete: %v", err)
		}
		// Verify deletion: Query should return gorm.ErrRecordNotFound
		// Always use errors.Is to check for gorm.ErrRecordNotFound
		err := db.First(&User{}, user.ID).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("expected not found, got %v", err)
		}
	})
}
