package main

import (
	. "github.com/otiai10/mint"
	"testing"
)

func TestCreateConfig(t *testing.T) {
	Expect(t, createDirConfig(""))
}
