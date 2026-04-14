package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	envOnce   sync.Once
	envLoaded bool
)

type DBType string

const (
	DBTypeSQLite   DBType = "sqlite"
	DBTypePostgres DBType = "postgres"
	DBTypeMySQL    DBType = "mysql"
)

/*
LoadEnv loads the .env file in the parent directory of the testutil package.
*/
func loadEnv() {
	envOnce.Do(func() {
		// Get directory of the current file, which is the testutil package
		_, currentFile, _, ok := runtime.Caller(0)
		if !ok {
			return
		}

		testutilDir := filepath.Dir(currentFile)
		gormExampleDir := filepath.Dir(testutilDir)
		envPath := filepath.Join(gormExampleDir, ".env")
		if err := godotenv.Load(envPath); err != nil {
			return
		}
		envLoaded = true
	})
}

func getDBType() DBType {
	loadEnv()
	dbType := os.Getenv("TEST_DB_TYPE")
	switch dbType {
	case "mysql":
		return DBTypeMySQL
	case "postgres", "postgresql":
		return DBTypePostgres
	default:
		return DBTypeSQLite
	}
}

func getDBDir() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil
	}
	testutilDir := filepath.Dir(currentFile)
	gormExampleDir := filepath.Dir(testutilDir)
	dbDir := filepath.Join(gormExampleDir, "db")
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return "", err
	}
	return dbDir, nil
}

func NewTestDB(t *testing.T, filename string) *gorm.DB {
	t.Helper()
	dbTyper := getDBType()
	var db *gorm.DB
	var err error
	switch dbTyper {
	case DBTypeMySQL:
		db, err = newMySQLDB(t)
	case DBTypePostgres:
		db, err = newPostgresDB(t)
	case DBTypeSQLite:
		db, err = newSQLiteDB(t, filename)
	default:
		t.Fatalf("unsupported DB type: %s", dbTyper)
	}
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	return db
}
func newSQLiteDB(t *testing.T, filename string) (*gorm.DB, error) {
	dbDir, err := getDBDir()
	if err != nil {
		return nil, err
	}
	if filename == "" {
		filename = "test.sqlite.db"
	} else {
		ext := filepath.Ext(filename)
		if ext == "" {
			ext = ".db"
		}
		base := filename[:len(filename)-len(ext)]
		if base == "" {
			base = "test"
		}
		baseLower := strings.ToLower(base)
		if !strings.Contains(baseLower, "sqlite") {
			filename = base + "_sqlite" + ext
		} else {
			filename = base + ext
		}
	}
	dbPath := filepath.Join(dbDir, filename)
	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // Prefix for all table names (e.g., "app_")
			SingularTable: false, // Use singular table names (User -> user instead of users)
			NoLowerCase:   false, // Disable automatic lowercasing
			NameReplacer:  nil,   // Custom name replacer function
		},
	})
}
func newMySQLDB(t *testing.T) (*gorm.DB, error) {
	loadEnv()
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
		t.Logf("using default MySQL DSN, set TEST_MYSQL_DSN in .env file or environment variable to override")
	}

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: Set to logger.Info to see all SQL queries in development
		Logger: logger.Default.LogMode(logger.Silent),

		// NamingStrategy: Customize table and column naming
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
			NoLowerCase:   false,
		},
	})
}
func newPostgresDB(t *testing.T) (*gorm.DB, error) {
	loadEnv()
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password dbname=testdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		t.Logf("using default PostgreSQL DSN, set TEST_POSTGRES_DSN in .env file or environment variable to override")
	}

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: Set to logger.Info to see all SQL queries in development
		Logger: logger.Default.LogMode(logger.Silent),

		// NamingStrategy: Customize table and column naming
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
			NoLowerCase:   false,
		},
	})
}
