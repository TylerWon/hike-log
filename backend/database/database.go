package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database is a thin wrapper around the gorm.DB type.
type Database interface {
	AutoMigrate(dst ...interface{}) error
	Create(value interface{}) (tx *gorm.DB)
	DB() (*sql.DB, error)
	Exec(sql string, values ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Order(value interface{}) (tx *gorm.DB)
	Preload(query string, args ...interface{}) (tx *gorm.DB)
	Raw(sql string, values ...interface{}) (tx *gorm.DB)
}

// databaseImpl is an implementation of the Database interface.
type databaseImpl struct {
	db *gorm.DB
}

// DBConfig stores database connection details.
type DBConfig struct {
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string
}

// Connects to the database specified in the given dbConfig.
func Connect(dbConfig DBConfig) (Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s",
		dbConfig.DbHost,
		dbConfig.DbPort,
		dbConfig.DbName,
		dbConfig.DbUser,
		dbConfig.DbPassword,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &databaseImpl{db}, nil
}

// Runs migration for given the models.
func (database *databaseImpl) AutoMigrate(dst ...interface{}) error {
	return database.db.AutoMigrate(dst...)
}

// Returns `*sql.DB`.
func (database *databaseImpl) DB() (*sql.DB, error) {
	return database.db.DB()
}

// Inserts value, returning the inserted data's primary key in value's id
func (database *databaseImpl) Create(value interface{}) (tx *gorm.DB) {
	return database.db.Create(value)
}

// Executes raw SQL DDL (INSERT, UPDATE, DELETE, etc.)
func (database *databaseImpl) Exec(sql string, values ...interface{}) (tx *gorm.DB) {
	return database.db.Exec(sql, values...)
}

// Finds all records matching given conditions conds
func (database *databaseImpl) Find(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return database.db.Find(dest, conds...)
}

// Finds the first record ordered by primary key, matching given conditions conds
func (database *databaseImpl) First(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return database.db.First(dest, conds...)
}

// Specifies the order when retrieving records from database
func (database *databaseImpl) Order(value interface{}) (tx *gorm.DB) {
	return database.db.Order(value)
}

// Preload associations with given conditions
func (database *databaseImpl) Preload(query string, args ...interface{}) (tx *gorm.DB) {
	return database.db.Preload(query, args...)
}

// Executes a raw SQL query (SELECT)
func (database *databaseImpl) Raw(sql string, values ...interface{}) (tx *gorm.DB) {
	return database.db.Raw(sql, values...)
}
