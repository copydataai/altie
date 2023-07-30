package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	c := require.New(t)

	homeDir, err := GetHomeDir()
	c.NoError(err)
	c.NotEmpty(homeDir)

	config, err := readTomlConfig(homeDir)
	c.NoError(err)
	c.Equal(*config, ConfigThemes{
		Config: Config{
			ThemesDirectory: fmt.Sprintf(RouteThemes, homeDir),
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{},
			LastMod:  time.Time{},
			FontSize: 14,
			Font:     "monoscape",
		},
	})

	config, err = readTomlConfig("/usr/share")
	c.Error(err)
	c.Nil(config)
}

func TestCheckLastModThemes(t *testing.T) {
	c := require.New(t)

	homeDir, err := GetHomeDir()
	c.NoError(err)
	c.NotEmpty(homeDir)

	timeNow := time.Now()

	isModified, err := checkLastModThemes(homeDir, timeNow)
	c.NoError(err)
	c.False(isModified)

	timeZero := time.Time{}

	isModified, err = checkLastModThemes(homeDir, timeZero)
	c.NoError(err)
	c.True(isModified)
}

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
