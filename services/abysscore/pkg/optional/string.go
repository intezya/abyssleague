package optional

import "time"

type (
	String = Optional[string]
	Time   = Optional[time.Time]
	Int    = Optional[int]
)
