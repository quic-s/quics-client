package defaultdb

type DB struct {
	Crud CRUD
}

func DBProvider(crud CRUD) (*DB, error) {
	return &DB{
		Crud: crud,
	}, nil
}

type CRUD interface {
	Update(key string, value []byte) error
	View(key string) ([]byte, error)
	Delete(key string) error
}
