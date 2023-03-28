package rungcal

import "os"

type Option struct {
	TargetDate string
	Project    string
	Verbose    bool
	DryRun     bool
}

func getEnv(key string, defaulValue ...string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaulValue[0]
}
