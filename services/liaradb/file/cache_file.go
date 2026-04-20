package file

type CacheFile struct {
	File
	name   string
	fsys   *Cache
	closed bool
}

func (cf *CacheFile) Close() error {
	return cf.fsys.CloseFile(cf.name)
}

func (cf *CacheFile) closeFile() error {
	err := cf.File.Close()
	if err != nil {
		return err
	}

	cf.closed = true
	return nil
}

func (cf *CacheFile) IsOpen() bool {
	return !cf.closed
}
