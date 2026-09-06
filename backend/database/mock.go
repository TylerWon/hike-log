package database

import (
	"database/sql"

	"gorm.io/gorm"
)

type MockDatabase struct {
	AutoMigrateError error
	CreateResult     *gorm.DB
	DBResult         *sql.DB
	DBError          error
	ExecResult       *gorm.DB
	FindResult       *gorm.DB
	FirstResult      *gorm.DB
	RawResult        *gorm.DB
}

func (m *MockDatabase) AutoMigrate(dst ...interface{}) error {
	return m.AutoMigrateError
}

func (m *MockDatabase) Create(value interface{}) (tx *gorm.DB) {
	return m.CreateResult
}

func (m *MockDatabase) DB() (*sql.DB, error) {
	return m.DBResult, m.DBError
}

func (m *MockDatabase) Exec(sql string, values ...interface{}) (tx *gorm.DB) {
	return m.ExecResult
}

func (m *MockDatabase) Find(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return m.FindResult
}

func (m *MockDatabase) First(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return m.FirstResult
}

func (m *MockDatabase) Order(value interface{}) Database {
	return m
}

func (m *MockDatabase) Preload(query string, args ...interface{}) Database {
	return m
}

func (m *MockDatabase) Raw(sql string, values ...interface{}) (tx *gorm.DB) {
	return m.RawResult
}
