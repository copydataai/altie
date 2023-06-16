package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	defaultFont     = "monoscape"
	defaultFontSize = 14
	// DirectoryThemes is to encapsulate the safe place of variables
	DirectoryThemes = "%s/themes"
	// RouteConfig is a const to replace by $HOME/.altie/altie.conf
	RouteConfig = "%s/.altie/altie.conf"
	// RouteThemes use $HOME/.altie/themes
	RouteThemes = "%s/.altie/themes"
)

type Config struct {
	HomeDirectory string `toml:"HomeDirectory"`
}

type ThemeConfig struct {
	FontSize int64  `toml:"FontSize"`
	Font     string `toml:"Font"`
}

type ConfigThemes struct {
	Config      `toml:"Config"`
	ThemeConfig `toml:"ConfigTheme"`
}

func createDirConfig(mainDirectory string) (string, error) {
	configDirectory := fmt.Sprintf(RouteConfig, mainDirectory)
	err := os.MkdirAll(mainDirectory+"/.altie", os.ModePerm)
	if err != nil {
		return "", err
	}

	return configDirectory, nil
}

func CreateConfig(mainDirectory string) error {
	configDirectory, err := createDirConfig(mainDirectory)
	if err != nil {
		return err
	}

	configFile, err := os.Create(configDirectory)
	if err != nil {
		return err
	}

	defer configFile.Close()

	defaultConfig := &ConfigThemes{
		Config{
			HomeDirectory: configDirectory,
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
