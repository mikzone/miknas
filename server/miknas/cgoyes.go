//go:build cgo

package miknas

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite3(dbpath string, gormConf gorm.Option) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dbpath), gormConf)
}
