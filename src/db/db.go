package db

import (
	"crud-example-go/src/errors"
	"crud-example-go/src/intreface/db"
	"gorm.io/gorm"
)

type DbWrapper struct {
	DB *gorm.DB
}

func (d *DbWrapper) Find(dest interface{}, conds ...interface{}) (tx db.DB) {
	return &DbWrapper{DB: d.DB.Find(dest, conds...)}
}

func (d *DbWrapper) First(dest interface{}, conds ...interface{}) (tx db.DB) {
	return &DbWrapper{DB: d.DB.First(dest, conds...)}
}

func (d *DbWrapper) Save(value interface{}) (tx db.DB) {
	return &DbWrapper{DB: d.DB.Save(value)}
}

func (d *DbWrapper) Delete(value interface{}, conds ...interface{}) (tx db.DB) {
	return &DbWrapper{DB: d.DB.Delete(value, conds...)}
}

func (d *DbWrapper) Preload(query string, args ...interface{}) (tx db.DB) {
	return &DbWrapper{DB: d.DB.Preload(query, args...)}
}

func (d *DbWrapper) GetError() error {
	err := d.DB.Error
	if err == gorm.ErrRecordNotFound {
		return errors.ErrRecordNotFound
	}
	return err
}
