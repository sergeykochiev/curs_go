package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb(db *gorm.DB, schema_file string) error {
	data, err := os.ReadFile(schema_file)
	if err != nil {
		log.Fatal("E cannot read schema file ("+schema_file+"): ", err.Error())
	}
	res := db.Exec(string(data))
	if res.Error != nil {
		println("E initializing db failed: ", res.Error.Error())
	}
	return err
}

func ConnectDb(filepath string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(filepath), &gorm.Config{})
}
