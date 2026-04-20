package file

type file struct {
	File
	name string
	fsys *FileSsytemCache
}

func (f *file) Close() error {
	return f.fsys.CloseFile(f.name)
}
