package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/BurntSushi/toml"
	"github.com/copydataai/altie/internal/config"
	"github.com/copydataai/altie/internal/themes"
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

	appConfig := config.NewAppConfig(tmpDir)

	err = os.MkdirAll(appConfig.AlacrittyDir, os.ModePerm)
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

	// Error when a altie.conf doesn't exists
	// But the directory exists
	err = CreateConfig()
	c.Error(err)
}

func TestListThemes(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err)

	defer os.RemoveAll(tmpDir)

	appConfig := config.NewAppConfig(tmpDir)

	err = os.MkdirAll(appConfig.ThemesDir, os.ModePerm)
	c.NoError(err)

	them := []string{"3024.dark", "3024.light", "Afterglow", "Argonaut"}
	for _, theme := range them {
		themeFileName := theme + ".yml"
		file, _ := os.Create(appConfig.ThemesDir + "/" + themeFileName)

		defer file.Close()
	}
	go func() {
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	configThemes := &config.ConfigThemes{
		Config: config.Config{
			ThemesDirectory: appConfig.ThemesDir,
		},
		ThemeConfig: config.ThemeConfig{
			Themes:   them,
			LastMod:  "",
			FontSize: 0,
			Font:     "",
		},
	}
	err = ListThemes(configThemes, appConfig)
	c.Error(err)

	go func() {
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	err = ListThemes(configThemes, appConfig)
	c.Error(err)
	c.True(os.IsNotExist(err))

	err = os.MkdirAll(appConfig.AlacrittyDir, os.ModePerm)
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

	fakeAppConfig := config.NewAppConfig("//")

	err = ListThemes(configThemes, fakeAppConfig)
	c.Error(err)
	fmt.Println(err.Error())
	c.True(os.IsNotExist(err))

	configThemes.Config.ThemesDirectory = appConfig.ThemesDir + "/fake"
	err = ListThemes(configThemes, appConfig)

	c.Error(err)
	c.EqualError(err, fmt.Errorf("no options provided").Error())

}

func TestApplyFontTheme(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err)

	defer os.RemoveAll(tmpDir)

	appConfig := config.NewAppConfig(tmpDir)

	err = os.MkdirAll(appConfig.ThemesDir, os.ModePerm)
	c.NoError(err)

	them := []string{"3024.dark", "3024.light", "Afterglow", "Argonaut"}
	for _, theme := range them {
		themeFileName := theme + ".yml"
		file, _ := os.Create(appConfig.ThemesDir + "/" + themeFileName)

		defer file.Close()
	}
	// TODO: Modify this config
	configThemes := &config.ConfigThemes{
		Config: config.Config{
			ThemesDirectory: appConfig.ThemesDir,
		},
		ThemeConfig: config.ThemeConfig{
			Themes:   them,
			LastMod:  "",
			FontSize: 12,
			Font:     "SpaceMono Nerd Font",
		},
	}

	err = os.MkdirAll(appConfig.AlacrittyDir, os.ModePerm)
	c.NoError(err)

	f, err := os.Create(appConfig.AlacrittyConfig)
	c.NoError(err)

	err = toml.NewEncoder(f).Encode(configThemes)
	c.NoError(err)

	c.NoError(f.Close())

	err = config.CreateConfig(appConfig)
	c.NoError(err)

	go func() {
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	err = ListThemes(configThemes, appConfig)
	c.NoError(err)

	f, err = os.Open(appConfig.ConfigDir)
	c.NoError(err)

	prevConfig, err := themes.CheckAlacrittyConfig(appConfig.AlacrittyConfig)
	c.NoError(err)

	fmt.Println(prevConfig)

	go func() {
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	fakeAppConfig := config.NewAppConfig("//")

	err = ListThemes(configThemes, fakeAppConfig)
	c.Error(err)
	fmt.Println(err.Error())
	c.True(os.IsNotExist(err))

	configThemes.Config.ThemesDirectory = appConfig.ThemesDir + "/fake"
	err = ListThemes(configThemes, appConfig)
	c.Error(err)
	c.EqualError(err, fmt.Errorf("no options provided").Error())
}
