package tablename

import "strings"

type Extension string

const (
	IndexExtension      Extension = "idx"
	KeyValueExtension   Extension = "kv"
	EventLogExtension   Extension = "el"
	OutboxExtension     Extension = "out"
	RequestLogExtension Extension = "rl"
	TenantExtension     Extension = "tn"
)

func (e Extension) String() string { return string(e) }

func IsIndex(fn string) bool {
	return strings.HasSuffix(fn, IndexExtension.String())
}
