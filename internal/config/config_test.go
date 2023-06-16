package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateConfig(t *testing.T) {
	c := require.New(t)

	homeDir, err := GetHomeDir()
	c.NoError(err)
	c.NotEmpty(homeDir)

	err = CreateConfig(homeDir)
	c.NoError(err)

	err = CreateConfig("/usr/share/")
	c.Error(err)

	err = os.Unsetenv("HOME")
	c.NoError(err)

	errDir, err := GetHomeDir()
	c.Error(err)
	c.Empty(errDir)
}
