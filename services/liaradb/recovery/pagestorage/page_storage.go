package pagestorage

type PageStorage struct {
}

func New() *PageStorage {
	return &PageStorage{}
}

func (t *PageStorage) Append([]byte) error {
	return nil
}

func (t *PageStorage) Sync([]byte) error {
	return nil
}
