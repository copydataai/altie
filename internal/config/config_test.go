package config

import (
	"os"
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

	appConfig := NewAppConfig(tmpDir)
	err = CreateConfig(appConfig)
	c.NoError(err, "Failed to create config: %v", err)

	// Check if the config directory exists
	_, err = os.Stat(appConfig.ConfigDir)
	c.False(os.IsNotExist(err), "Config directory does not exist")

	// Check if the config file exists
	_, err = os.Stat(appConfig.ConfigFilePath)
	c.False(os.IsNotExist(err), "Config file does not exist")

	// Check if the default config is correct
	expectedConfig := &ConfigThemes{
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
	actualConfig, err := CheckConfig(appConfig.ConfigFilePath)
	c.NoError(err, "Failed to read config: %v", err)

	c.Equal(*expectedConfig, *actualConfig)
}

func TestReadConfig(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err, "Failed to create temporary directory: %v", err)

	defer os.RemoveAll(tmpDir)

	appConfig := NewAppConfig(tmpDir)

	err = CreateConfig(appConfig)
	c.NoError(err, "Failed to create config: %v", err)

	config, err := CheckConfig(appConfig.ConfigFilePath)
	c.NoError(err)
	c.Equal(*config, ConfigThemes{
		Config: Config{
			ThemesDirectory: appConfig.ThemesDir,
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: defaultFontSize,
			Font:     defaultFont,
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

	appConfig := NewAppConfig(tmpDir)

	err = CreateConfig(appConfig)
	c.NoError(err)

	config, err := CheckConfig(appConfig.ConfigFilePath)
	c.NoError(err)
	c.EqualValues(&ConfigThemes{
		Config: Config{
			ThemesDirectory: appConfig.ThemesDir,
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: defaultFontSize,
			Font:     defaultFont,
		},
	}, config)

	timeNow := time.Now()

	err = config.SetModifiedThemes(appConfig, timeNow, []string{"Hello", "world", "again"})
	c.NoError(err)

	config, err = CheckConfig(appConfig.ConfigFilePath)
	c.NoError(err)
	c.EqualValues(&ConfigThemes{
		Config: Config{
			ThemesDirectory: appConfig.ThemesDir,
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

	appConfig := NewAppConfig(tmpDir)

	err = CreateConfig(appConfig)
	c.NoError(err, "Failed to create config: %v", err)

	config, err := CheckConfig(appConfig.ConfigFilePath)
	c.NoError(err)
	c.NotNil(config)
	c.EqualValues(&ConfigThemes{
		Config: Config{
			ThemesDirectory: appConfig.ThemesDir,
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: defaultFontSize,
			Font:     defaultFont,
		},
	}, config)

	confFile, err := os.OpenFile(appConfig.ConfigFilePath, os.O_RDWR, os.ModePerm)
	c.NoError(err)

	config.ThemeConfig.Font = "Mononoki"
	config.ThemeConfig.FontSize = 64

	err = encodeTomlConfig(confFile, config)
	c.NoError(err)

	confFile.Close()

	config, err = CheckConfig(appConfig.ConfigFilePath)
	c.NoError(err)
	c.EqualValues(&ConfigThemes{
		Config: Config{
			ThemesDirectory: appConfig.ThemesDir,
		},
		ThemeConfig: ThemeConfig{
			Themes:   []string{},
			LastMod:  "",
			FontSize: 64,
			Font:     "Mononoki",
		},
	}, config)

	confFile, err = os.Open(appConfig.ConfigDir)
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
}

func TestCheckLastModThemes(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err, "Failed to create temporary directory: %v", err)

	appConfig := NewAppConfig(tmpDir)

	err = os.MkdirAll(appConfig.ThemesDir, os.ModePerm)
	c.NoError(err)

	defer os.RemoveAll(tmpDir)

	timeNow := time.Now()

	isModified, err := checkLastModThemes(appConfig.HomeDir, timeNow)
	c.NoError(err)
	c.False(isModified)

	timeZero := time.Time{}

	isModified, err = checkLastModThemes(appConfig.HomeDir, timeZero)
	c.NoError(err)
	c.True(isModified)
}
