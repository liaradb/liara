package tablename

import (
	"fmt"
	"strings"

	"github.com/liaradb/liaradb/storage/link"
)

type Extension string

const (
	CollectionExtension Extension = ".col"
	BTreeExtension      Extension = ".bt"
)

func (e Extension) String() string { return string(e) }

func (e Extension) FileName(n string) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v%v", n, e))
}

func IsBTree(fn string) bool {
	return strings.HasSuffix(fn, BTreeExtension.String())
}

func IsCollection(fn string) bool {
	return strings.HasSuffix(fn, CollectionExtension.String())
}
