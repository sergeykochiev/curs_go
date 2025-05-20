package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ExecuteFile(db *gorm.DB, file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	res := db.Exec(string(data))
	return res.Error
}

func Connect(filepath string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(filepath), &gorm.Config{})
}
