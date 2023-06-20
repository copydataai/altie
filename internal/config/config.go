package config

import (
	"fmt"
	"os"

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
	// AlacrittyConfigDir is $HOME/.config/alacritty/alacritty.yml
	AlacrittyConfigDir = "%s/.config/alacritty/alacritty.yml"
)

// TODO: Implement a method to read ThemesDirectory
type Config struct {
	ThemesDirectory string `toml:"ThemesDirectory"`
}

// TODO: Implement a method to read and don't modify the themes
type ThemeConfig struct {
	FontSize int64  `toml:"FontSize"`
	Font     string `toml:"Font"`
}

type ConfigThemes struct {
	Config      `toml:"Config"`
	ThemeConfig `toml:"ConfigTheme"`
}

func createDirConfig(homeDir string) error {
	configDir := fmt.Sprintf(ConfigDir, homeDir)
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func CreateConfig(homeDir string) error {
	err := createDirConfig(homeDir)
	if err != nil {
		return err
	}

	configFile, err := os.Create(fmt.Sprintf(RouteConfig, homeDir))
	if err != nil {
		return err
	}

	defer configFile.Close()

	defaultConfig := &ConfigThemes{
		Config{
			ThemesDirectory: fmt.Sprintf(RouteThemes, homeDir),
		},
		ThemeConfig{
			FontSize: defaultFontSize,
			Font:     defaultFont,
		},
	}

	err = createTomlConfig(configFile, defaultConfig)
	if err != nil {
		return err
	}

	return nil
}

func createTomlConfig(configFile *os.File, configTheme *ConfigThemes) error {
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

func CheckConfig(homeDir string) (*ConfigThemes, error) {
	configAltie := &ConfigThemes{}
	configAltiePath := fmt.Sprintf(RouteConfig, homeDir)

	_, err := toml.DecodeFile(configAltiePath, configAltie)
	if err != nil {
		return nil, err
	}

	return configAltie, nil
}
