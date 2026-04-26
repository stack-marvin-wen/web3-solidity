package basic

import (
	"L2-gorm/testutil"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestQuery(t *testing.T) {
	db := testutil.NewTestDB(t, "query.db")
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
	t.Run("query/basic", func(t *testing.T) {
		var users []User
		if err := db.Where("status = ?", "active").Find(&users).Error; err != nil {
			t.Fatalf("query user err: %v", err)
		}
		t.Logf("basic query: %v\n", users)
	})
	t.Run("query/mutiplecondition", func(t *testing.T) {
		var users []User
		if err := db.Where("status = ? and age > ?", "active", 25).Find(&users).Error; err != nil {
			t.Fatalf("query user err: %v", err)
		}
		t.Logf("mutiplecondition query: %v\n", users)
	})
	t.Run("query/like", func(t *testing.T) {
		var users []User
		if err := db.Where("name like ?", "%j%").Find(&users).Error; err != nil {
			t.Fatalf("query user err: %v", err)
		}
		t.Logf("query like: %v\n", users)

	})
	t.Run("query/in", func(t *testing.T) {
		var users []User
		if err := db.Where("name in ?", []string{"Alice", "Bob", "Charlie"}).Find(&users).Error; err != nil {
			t.Fatalf("query user err: %v", err)
		}
		t.Logf("query in: %v\n", users)
	})
	t.Run("query/between", func(t *testing.T) {
		var users []User
		if err := db.Where("age BETWEEN ? AND ?", 20, 30).Find(&users).Error; err != nil {
			t.Fatalf("query user err: %v", err)
		}
		t.Logf("query between: %v\n", users)
	})

	t.Run("query/order", func(t *testing.T) {
		var users []User
		if err := db.Order("name desc").Find(&users).Error; err != nil {
			t.Fatalf("query user err: %v", err)
		}
		t.Logf("query order: %v\n", users)
	})
	t.Run("query/limit_offset_page", func(t *testing.T) {
		pageSize := 5
		page := 2
		offset := (page - 1) * pageSize
		var users []User
		if err := db.Order("name desc").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
			t.Fatalf("query user err: %v", err)
		}
		t.Logf("query limit offset page: %v\n", users)
	})
	t.Run("query/scopes_active", func(t *testing.T) {
		var users []User
		if err := db.Scopes(activeUser()).Find(&users).Error; err != nil {
			t.Fatalf("query user err: %v", err)
		}
		t.Logf("query scopes active: %v\n", users)
	})
	t.Run("query/paginate", func(t *testing.T) {
		var users []User
		if err := db.Scopes(paginate(1, 10)).Find(&users).Error; err != nil {
			t.Fatalf("query user err: %v", err)
		}
		t.Logf("query scopes active: %v\n", users)
	})
}

func activeUser() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "active")
	}
}
func paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset(pageSize * (page - 1))
	}
}
