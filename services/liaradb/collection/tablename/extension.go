package tablename

import "strings"

type Extension string

const (
	CollectionExtension Extension = ".col"
	BTreeExtension      Extension = ".bt"
)

func (e Extension) String() string { return string(e) }

func IsBTree(fn string) bool {
	return strings.HasSuffix(fn, BTreeExtension.String())
}

func IsCollection(fn string) bool {
	return strings.HasSuffix(fn, CollectionExtension.String())
}
