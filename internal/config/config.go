package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	defaultFont     = "monoscape"
	defaultFontSize = 14
	// TODO: Delete after implement check in Config
	ThemesDir = "%s/themes"
	// ConfigDir is $HOME/.altie directory
	ConfigDir = "%s/.altie"
	// RouteConfig is a const to replace by $HOME/.altie/altie.conf
	RouteConfig = "%s/.altie/altie.conf"
	// RouteThemes use $HOME/.altie/themes
	RouteThemes = "%s/.altie/themes"
	// AlacrittyDir is $HOME/.config/alacritty
	AlacrittyDir = "%s/.config/alacritty"
	// AlacrittyConfigDir is $HOME/.config/alacritty/alacritty.toml
	AlacrittyConfigDir = "%s/.config/alacritty/alacritty.toml"
)

type Config struct {
	ThemesDirectory string `toml:"ThemesDirectory"`
}

// TODO: Implement a method to read and don't modify the themes
type ThemeConfig struct {
	Themes   []string `toml:"Themes"`
	LastMod  string   `toml:"LastModified"`
	FontSize int64    `toml:"FontSize"`
	Font     string   `toml:"Font"`
}

type ConfigThemes struct {
	Config      `toml:"Config"`
	ThemeConfig `toml:"ConfigTheme"`
}

func checkLastModThemes(homeDir string, lastMod time.Time) (bool, error) {
	themesDir := fmt.Sprintf(RouteThemes, homeDir)

	info, err := os.Stat(themesDir)
	if err != nil {
		return false, err
	}

	difference := info.ModTime().Compare(lastMod)
	if difference < 1 {
		return false, nil
	}

	return true, nil
}

func (config *ConfigThemes) SetModifiedThemes(homeDir string, lastMod time.Time, listThemes []string) error {
	config.LastMod = lastMod.Format(time.RFC3339)
	config.ThemeConfig.Themes = listThemes

	configFile, err := os.Create(fmt.Sprintf(RouteConfig, homeDir))
	if err != nil {
		return err
	}

	defer configFile.Close()

	err = encodeTomlConfig(configFile, config)
	if err != nil {
		return err
	}

	return nil
}

func CreateConfig(mainDir string) error {
	configDir := fmt.Sprintf(ConfigDir, mainDir)

	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	configFile, err := os.Create(fmt.Sprintf(RouteConfig, mainDir))
	if err != nil {
		return err
	}

	defer configFile.Close()

	defaultConfig := &ConfigThemes{
		Config{
			ThemesDirectory: fmt.Sprintf(RouteThemes, mainDir),
		},
		ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: defaultFontSize,
			Font:     defaultFont,
		},
	}

	err = encodeTomlConfig(configFile, defaultConfig)
	if err != nil {
		return err
	}

	return nil
}

func encodeTomlConfig(configFile *os.File, configTheme *ConfigThemes) error {
	err := toml.NewEncoder(configFile).Encode(configTheme)
	if err != nil {
		return err
	}

	return nil
}

func GetHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir, nil
}

func CheckConfig(configDir string) (*ConfigThemes, error) {
	configAltie := &ConfigThemes{}

	_, err := toml.DecodeFile(configDir, configAltie)
	if err != nil {
		return nil, err
	}

	return configAltie, nil
}
