package record

import "testing"

func TestCollection_String(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		collection Collection
		result     string
	}{
		"should handle zero": {
			collection: 0,
			result:     ""},
		"should handle system": {
			collection: CollectionSystem,
			result:     "system"},
		"should handle request": {
			collection: CollectionRequest,
			result:     "request"},
		"should handle outbox": {
			collection: CollectionOutbox,
			result:     "outbox"},
		"should handle event": {
			collection: CollectionEvent,
			result:     "event"},
		"should handle value": {
			collection: CollectionValue,
			result:     "value"},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if r := c.collection.String(); r != c.result {
				t.Errorf("incorrect string: %v, expected: %v", r, c.result)
			}
		})
	}
}
