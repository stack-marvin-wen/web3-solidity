package basic

import (
	"L2-gorm/testutil"
	"testing"
	"time"
)

func TestSetup(t *testing.T) {
	db := testutil.NewTestDB(t, "setup.db")
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
		t.Fatal(err)
	}
	if err := db.Exec("delete from users").Error; err != nil {
		t.Fatal(err)
	}
	users := []User{
		{Name: "Alice", Email: "alice@example.com", Age: 18, Status: "active"},
		{Name: "Bob", Email: "bob@example.com", Age: 19, Status: "active"},
		{Name: "Charlie", Email: "charlie@example.com", Age: 20, Status: "inactive"},
		{Name: "David", Email: "david@example.com", Age: 21, Status: "active"},
	}
	if err := db.CreateInBatches(users, 4).Error; err != nil {
		t.Fatalf("seed users: %v", err)
	}

	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		t.Fatalf("count users: %v", err)
	}
	if count != int64(len(users)) {
		t.Errorf("expected 4 users, got %d", count)
	}
	t.Logf("created %d users", count)
}
