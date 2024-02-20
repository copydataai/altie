package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	defaultFont     = "monoscape"
	defaultFontSize = 14
)

// AppConfig holds all application configuration paths
type AppConfig struct {
	HomeDir         string
	ConfigDir       string
	ConfigFilePath  string
	ThemesDir       string
	AlacrittyDir    string
	AlacrittyConfig string
}

func NewAppConfig(homeDir string) *AppConfig {
	baseDir := filepath.Join(homeDir, ".altie")
	return &AppConfig{
		HomeDir:         homeDir,
		ConfigDir:       baseDir,
		ConfigFilePath:  filepath.Join(baseDir, "altie.conf"),
		ThemesDir:       filepath.Join(baseDir, "themes"),
		AlacrittyDir:    filepath.Join(homeDir, ".config", "alacritty"),
		AlacrittyConfig: filepath.Join(homeDir, ".config", "alacritty", "alacritty.toml"),
	}
}

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

func checkLastModThemes(themesDir string, lastMod time.Time) (bool, error) {
	info, err := os.Stat(themesDir)
	if err != nil {
		return false, err
	}

	return info.ModTime().After(lastMod), nil
}

func (config *ConfigThemes) SetModifiedThemes(appConfig *AppConfig, lastMod time.Time, listThemes []string) error {
	config.LastMod = lastMod.Format(time.RFC3339)
	config.ThemeConfig.Themes = listThemes

	configFile, err := os.Create(appConfig.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	defer configFile.Close()

	return encodeTomlConfig(configFile, config)
}

func CreateConfig(appConfig *AppConfig) error {
	err := os.MkdirAll(appConfig.ConfigDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile, err := os.Create(appConfig.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	defer configFile.Close()

	defaultConfig := &ConfigThemes{
		Config{
			ThemesDirectory: appConfig.ThemesDir,
		},
		ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: defaultFontSize,
			Font:     defaultFont,
		},
	}

	return encodeTomlConfig(configFile, defaultConfig)
}

func encodeTomlConfig(configFile *os.File, configTheme *ConfigThemes) error {
	err := toml.NewEncoder(configFile).Encode(configTheme)
	if err != nil {
		return fmt.Errorf("failed to encode TOML config: %w", err)
	}

	return nil
}

func GetHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return homeDir, nil
}

func CheckConfig(configFilePath string) (*ConfigThemes, error) {
	configAltie := &ConfigThemes{}
	if _, err := toml.DecodeFile(configFilePath, configAltie); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return configAltie, nil
}
