package essql

type (
	Rows interface {
		Row
		Next() bool
	}

	Row interface {
		Scan(dest ...any) error
	}
)
