package storage

type DBStorage struct {
	filePath string
}

func NewDBStorage() *DBStorage {
	return &DBStorage{}
}

func (s *DBStorage) Update() interface{} {
	return nil
}

func (s *DBStorage) GetAll() map[string]interface{} {
	return nil
}
