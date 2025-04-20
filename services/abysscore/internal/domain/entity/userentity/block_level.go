package userentity

type SearchBlockLevel int
type AccountBlockLevel int

const (
	SearchBlockLevelWarning1 SearchBlockLevel = iota
	SearchBlockLevelWarning2
	SearchBlockLevelBan1h
	SearchBlockLevelBan6h
	SearchBlockLevelBan12h
	SearchBlockLevelBan24h
	SearchBlockLevelBan72h
	SearchBlockLevelBanInf = -1
)

const (
	AccountBlockLevelWarning AccountBlockLevel = iota
	AccountBlockLevelBan6h
	AccountBlockLevelBan24h
	AccountBlockLevelBan72h
	AccountBlockLevelBanInf = -1
)
