package liara

type (
	AggregateName string
	EventName     string
	GlobalVersion int
	Limit         int
	Schema        string
	Version       int
)

func (an AggregateName) String() string { return string(an) }
func (n EventName) String() string      { return string(n) }
func (gv GlobalVersion) Value() int     { return int(gv) }
func (l Limit) Value() int              { return int(l) }
func (s Schema) String() string         { return string(s) }
func (v Version) Value() int            { return int(v) }
