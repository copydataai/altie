package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/copydataai/altie/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCreateConfig(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err, "Failed to create temporary directory: %v", err)

	err = os.Unsetenv("HOME")
	c.NoError(err)

	err = CreateConfig()
	c.Error(err)
	c.EqualError(err, "$HOME is not defined")

	err = os.Setenv("HOME", tmpDir)
	c.NoError(err)

	go func() {
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	err = CreateConfig()
	c.NoError(err)

	go func() {
		keyboard.SimulateKeyPress("Y")
		time.Sleep(2 * time.Microsecond)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	err = CreateConfig()
	c.NoError(err)

	err = os.RemoveAll(tmpDir)
	c.NoError(err)

	// TODO: catch errors in config.CreateConfig with the homeDir
	//
	//

	// TODO: catch second CheckConfig removing the config file
	//

	tmpDir, err = os.MkdirTemp("", "test")
	c.NoError(err, "Failed to create temporary directory: %v", err)

	defer os.RemoveAll(tmpDir)

	err = os.Setenv("HOME", tmpDir)
	c.NoError(err)

	configDir := fmt.Sprintf(config.AlacrittyConfigDir, tmpDir)

	err = os.MkdirAll(configDir, os.ModePerm)
	c.NoError(err)

	go func() {
		keyboard.SimulateKeyPress("Y")
		time.Sleep(2 * time.Microsecond)
		keyboard.SimulateKeyPress("Y")
		time.Sleep(2 * time.Microsecond)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	err = CreateConfig()
	c.NoError(err)
}

func TestListThemes(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err)

	defer os.RemoveAll(tmpDir)

	themesDir := fmt.Sprintf(config.RouteThemes, tmpDir)
	err = os.MkdirAll(themesDir, os.ModePerm)
	c.NoError(err)

	themes := []string{"3024.dark", "3024.light", "Afterglow", "Argonaut"}
	for _, theme := range themes {
		themeFileName := theme + ".yml"
		file, _ := os.Create(themesDir + "/" + themeFileName)

		defer file.Close()
	}
	go func() {
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	configThemes := &config.ConfigThemes{
		Config: config.Config{
			ThemesDirectory: themesDir,
		},
		ThemeConfig: config.ThemeConfig{
			Themes:   themes,
			LastMod:  "",
			FontSize: 0,
			Font:     "",
		},
	}
	err = ListThemes(configThemes, tmpDir)
	c.Error(err)

	go func() {
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	err = ListThemes(configThemes, tmpDir)
	c.Error(err)
	c.True(os.IsNotExist(err))

	configDir := fmt.Sprintf(config.AlacrittyConfigDir, tmpDir)

	err = os.MkdirAll(configDir, os.ModePerm)
	c.NoError(err)

	go func() {
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	go func() {
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	err = ListThemes(configThemes, "//")
	c.Error(err)
	fmt.Println(err.Error())
	c.True(os.IsNotExist(err))

	configThemes.Config.ThemesDirectory = themesDir + "/fake"
	err = ListThemes(configThemes, tmpDir)

	c.Error(err)
	c.EqualError(err, fmt.Errorf("no options provided").Error())

}
