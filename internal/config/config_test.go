package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateConfig(t *testing.T) {
	c := require.New(t)
	// Test case 1: Test creating config directory
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err, "Failed to create temporary directory: %v", err)

	defer os.RemoveAll(tmpDir)

	err = CreateConfig(tmpDir)
	c.NoError(err, "Failed to create config: %v", err)

	// Check if the config directory exists
	configDir := filepath.Join(tmpDir, ".altie")
	_, err = os.Stat(configDir)
	c.False(os.IsNotExist(err), "Config directory does not exist")

	// Check if the config file exists
	configFile := filepath.Join(configDir, "altie.conf")

	_, err = os.Stat(configFile)
	c.False(os.IsNotExist(err), "Config file does not exist")

	// Check if the default config is correct
	expectedConfig := &ConfigThemes{
		Config{
			ThemesDirectory: filepath.Join(configDir, "themes"),
		},
		ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: defaultFontSize,
			Font:     defaultFont,
		},
	}
	actualConfig, err := CheckConfig(configFile)
	c.NoError(err, "Failed to read config: %v", err)

	c.Equal(*expectedConfig, *actualConfig)
}
func TestReadConfig(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err, "Failed to create temporary directory: %v", err)

	defer os.RemoveAll(tmpDir)

	err = CreateConfig(tmpDir)
	c.NoError(err, "Failed to create config: %v", err)

	configDir := filepath.Join(tmpDir, ".altie", "altie.conf")

	config, err := CheckConfig(configDir)
	c.NoError(err)
	c.Equal(*config, ConfigThemes{
		Config: Config{
			ThemesDirectory: fmt.Sprintf(RouteThemes, tmpDir),
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: 14,
			Font:     "monoscape",
		},
	})

	config, err = CheckConfig("/usr/share")
	c.Error(err)
	c.Nil(config)
}

func TestSetModifiedThemes(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err, "Failed to create temporary directory: %v", err)

	defer os.RemoveAll(tmpDir)

	err = CreateConfig(tmpDir)
	c.NoError(err, "Failed to create config: %v", err)

	configFile := fmt.Sprintf(RouteConfig, tmpDir)

	config, err := CheckConfig(configFile)
	c.NoError(err)
	c.EqualValues(&ConfigThemes{
		Config: Config{
			ThemesDirectory: fmt.Sprintf(RouteThemes, tmpDir),
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: defaultFontSize,
			Font:     defaultFont,
		},
	}, config)

	timeNow := time.Now()

	err = config.SetModifiedThemes(tmpDir, timeNow, []string{"Hello", "world", "again"})
	c.NoError(err)

	config, err = CheckConfig(configFile)
	c.NoError(err)
	c.EqualValues(&ConfigThemes{
		Config: Config{
			ThemesDirectory: fmt.Sprintf(RouteThemes, tmpDir),
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{"Hello", "world", "again"},
			LastMod:  timeNow.Format(time.RFC3339),
			FontSize: defaultFontSize,
			Font:     defaultFont,
		},
	}, config)

}

func TestEncodeTomlConfig(t *testing.T) {
	// Test case 1: Decoding valid TOML file
	c := require.New(t)
	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err, "Failed to create temporary directory: %v", err)

	defer os.RemoveAll(tmpDir)

	err = CreateConfig(tmpDir)
	c.NoError(err, "Failed to create config: %v", err)

	configDir := filepath.Join(tmpDir, ".altie", "altie.conf")

	config, err := CheckConfig(configDir)
	c.NoError(err)
	c.NotNil(config)
	c.EqualValues(&ConfigThemes{
		Config: Config{
			ThemesDirectory: fmt.Sprintf(RouteThemes, tmpDir),
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: defaultFontSize,
			Font:     defaultFont,
		},
	}, config)

	confFile, err := os.OpenFile(configDir, os.O_RDWR, os.ModePerm)
	c.NoError(err)

	config.ThemeConfig.Font = "Mononoki"
	config.ThemeConfig.FontSize = 64

	err = encodeTomlConfig(confFile, config)
	c.NoError(err)

	confFile.Close()

	config, err = CheckConfig(configDir)
	c.NoError(err)
	c.EqualValues(&ConfigThemes{
		Config: Config{
			ThemesDirectory: fmt.Sprintf(RouteThemes, tmpDir),
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: 64,
			Font:     "Mononoki",
		},
	}, config)

	confFile, err = os.Open(configDir)
	c.NoError(err)

	confFile.Close()

	err = encodeTomlConfig(confFile, config)
	c.Error(err)
}

func TestGetHomeDir(t *testing.T) {
	c := require.New(t)
	// Test case 1: HomeDir is not empty
	homeDir, err := GetHomeDir()
	c.NoError(err)
	c.NotEmpty(homeDir)

	// Test case 2: Error when getting HomeDir
	os.Setenv("HOME", "") // Simulate empty HOME environment variable
	_, err = GetHomeDir()
	c.Error(err)
	c.EqualError(err, "$HOME is not defined")
}

func TestCheckLastModThemes(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err, "Failed to create temporary directory: %v", err)

	err = os.MkdirAll(fmt.Sprintf(RouteThemes, tmpDir), os.ModePerm)
	c.NoError(err)

	defer os.RemoveAll(tmpDir)

	timeNow := time.Now()

	isModified, err := checkLastModThemes(tmpDir, timeNow)
	c.NoError(err)
	c.False(isModified)

	timeZero := time.Time{}

	isModified, err = checkLastModThemes(tmpDir, timeZero)
	c.NoError(err)
	c.True(isModified)
}
