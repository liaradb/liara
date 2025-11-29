package btree

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

// TODO: Test this
type KeyEntry struct {
	key   Key
	block BlockPosition
}

func newKeyEntry(key Key, block BlockPosition) KeyEntry {
	return KeyEntry{
		key:   key,
		block: block,
	}
}

func (le KeyEntry) Size() int { return le.key.Size() + BlockPositionSize }

func (le KeyEntry) Write(w io.Writer) error {
	return raw.WriteAll(w,
		le.key,
		le.block,
	)
}

func (le *KeyEntry) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&le.key,
		&le.block)
}
