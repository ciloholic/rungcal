package rungcal

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Unsetenv("TEST_ENV")
	if getEnv("TEST_ENV", "default") != "default" {
		t.Error(`ng`)
	}

	os.Setenv("TEST_ENV", "ok")
	if getEnv("TEST_ENV", "default") != "ok" {
		t.Error(`ng`)
	}
}
