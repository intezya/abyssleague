package config

type PathConfig struct {
	Authentication struct {
		Self     string
		Register string
		Login    string
	}
	Other struct {
		Liveness   string
		Readliness string
		Pprof      string
	}
}

var pathConfig = &PathConfig{
	Authentication: struct {
		Self     string
		Register string
		Login    string
	}{
		Self:     "/auth",
		Register: "/register",
		Login:    "/login",
	},
	Other: struct {
		Liveness   string
		Readliness string
		Pprof      string
	}{Liveness: "", Readliness: "", Pprof: ""},
}

var notInfoLogging []string

// PathForLevelInfo if path is in NotInfoLogging return false and request must be logged as debug
func (c *Config) PathForLevelInfo(path string) bool {
	for _, v := range notInfoLogging {
		if path == v {
			return false
		}
	}
	return true
}
