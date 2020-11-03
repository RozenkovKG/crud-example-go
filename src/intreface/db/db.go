package db

type DB interface {
	Find(dest interface{}, conds ...interface{}) (tx DB)
	First(dest interface{}, conds ...interface{}) (tx DB)
	Save(value interface{}) (tx DB)
	Delete(value interface{}, conds ...interface{}) (tx DB)
	Preload(query string, args ...interface{}) (tx DB)
	GetError() error
}
