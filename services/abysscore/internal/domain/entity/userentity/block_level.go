package userentity

import "time"

// SearchBlockLevel represents search functionality restriction level
type SearchBlockLevel int

// AccountBlockLevel represents account-wide restriction level
type AccountBlockLevel int

const (
	SearchBlockLevelWarning1 SearchBlockLevel = iota
	SearchBlockLevelWarning2
	SearchBlockLevelBan1h
	SearchBlockLevelBan6h
	SearchBlockLevelBan12h
	SearchBlockLevelBan24h
	SearchBlockLevelBan72h

	SearchBlockDecrementTime = time.Hour * 24 * 3 // 3 days
)

const (
	AccountBlockLevelWarning AccountBlockLevel = iota
	AccountBlockLevelBan6h
	AccountBlockLevelBan24h
	AccountBlockLevelBan72h

	AccountBlockDecrementTime = time.Hour * 24 * 3 // 3 days
)
