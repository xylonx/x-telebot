package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/hints"
)

func NewSqliteConn(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db = db.Clauses(hints.New("MAX_EXECUTION_TIME = 2000")) // unit is ms
	return db, nil
}
