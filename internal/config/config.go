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
	// AlacrittyConfigDir is $HOME/.config/alacritty/alacritty.yml
	AlacrittyConfigDir = "%s/.config/alacritty/alacritty.yml"
)

type Config struct {
	ThemesDirectory string `toml:"ThemesDirectory"`
}

// TODO: Implement a method to read and don't modify the themes
type ThemeConfig struct {
	Themes   []string  `toml:"Themes"`
	LastMod  time.Time `toml:"LastModified"`
	FontSize int64     `toml:"FontSize"`
	Font     string    `toml:"Font"`
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

func saveTomlConfig(configFile string, config *ConfigThemes) error {
	_, err := toml.DecodeFile(configFile, config)
	if err != nil {
		return err
	}

	return nil
}

func (config *ConfigThemes) SetModifiedThemes(homeDir string, lastMod time.Time, listThemes []string) error {
	config.LastMod = lastMod
	config.ThemeConfig.Themes = listThemes
	err := saveTomlConfig(homeDir, config)
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
			LastMod:  time.Time{},
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

func CheckConfig(configDir string) (*ConfigThemes, error) {
	configAltie := &ConfigThemes{}

	_, err := toml.DecodeFile(configDir, configAltie)
	if err != nil {
		return nil, err
	}

	return configAltie, nil
}
