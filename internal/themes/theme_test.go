package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/copydataai/altie/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCreateThemes(t *testing.T) {
	c := require.New(t)

	homeDir, err := config.GetHomeDir()
	c.NoError(err)

	err = os.Setenv("PWD", "/home/copydataai/Documents/Personal/OpenSource/alacritty-themes-go-version")
	c.NoError(err)

	repoDirectory, err := GetRepoDirectory()
	c.NoError(err)

	repoDirectory = filepath.Dir(filepath.Dir(repoDirectory))

	themesDirectory := fmt.Sprintf(config.ThemesDir, repoDirectory)
	configDirectory := fmt.Sprintf(config.RouteThemes, homeDir)

	err = CreateThemes(configDirectory, themesDirectory)
	c.NoError(err)

	newDirectory, err := GetRepoDirectory()
	c.NoError(err)

	themesDirectory = fmt.Sprintf(config.ThemesDir, newDirectory)

	c.NoError(err)
	err = CreateThemes(configDirectory, themesDirectory)
	c.Error(err)
}

func TestApplyTheme(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err)

	defer os.RemoveAll(tmpDir)

	stackThemesDir, err := GetRepoDirectory()
	c.NoError(err)

	themesDir := filepath.Join(tmpDir, "themes")

	err = ApplyTheme(stackThemesDir, themesDir)
	c.NoError(err)

	// Test case 2: Copying theme file with invalid path
	err = ApplyTheme("", themesDir)
	c.Error(err)
}

func TestListThemes(t *testing.T) {
	// Test case 1: Empty directory
	c := require.New(t)
	emptyDir := ""
	expectedEmptyDir := []string{}

	resultEmptyDir, errEmptyDir := ListThemes(emptyDir)
	c.NoError(errEmptyDir)
	c.Equal(expectedEmptyDir, resultEmptyDir)

	// Test case 2: Directory with themes
	dirWithThemes, err := GetRepoDirectory()
	c.NoError(err)

	dirWithThemes = filepath.Dir(dirWithThemes)
	dirWithThemes = filepath.Dir(dirWithThemes)

	dirWithThemes = filepath.Join(dirWithThemes, "themes")

	expectedDirWithThemes := "3024.dark.yml"
	resultDirWithThemes, errDirWithThemes := ListThemes(dirWithThemes)
	c.NoError(errDirWithThemes)
	c.Contains(resultDirWithThemes, expectedDirWithThemes)
}

func TestBackUpTheme(t *testing.T) {
	// Test case 1: Verify that the backup file is created with the correct name
	// Set up test data
	c := require.New(t)
	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err)

	defer os.RemoveAll(tmpDir)

	// create a file with name alacritty.yml inside tmpDir
	alacrittyConfDir := fmt.Sprintf("%s/alacritty.yml", tmpDir)
	_, err = os.Create(alacrittyConfDir)
	c.NoError(err)

	// Call the function
	backupPath, err := BackUpTheme(alacrittyConfDir)
	c.NoError(err)

	// Verify the backup file name
	expectedFileName := fmt.Sprintf("%s.%d%d%d.bak", alacrittyConfDir, time.Now().Year(), time.Now().Month(), time.Now().Day())
	c.Equal(backupPath, expectedFileName)

	_, err = os.Stat(alacrittyConfDir)
	c.False(!os.IsNotExist(err))

	// Verify the backup file path exists
	_, err = os.Stat(backupPath)
	c.False(os.IsNotExist(err))

	alacrittyConfDir = "/path/to/alacritty/conf/dir"

	// Call the function with a non-existent directory
	backupPath, err = BackUpTheme(alacrittyConfDir)
	c.Error(err)
	c.Empty(backupPath)
}

func TestCheckAltieThemes(t *testing.T) {
	// Test case 1: When the directory exists
	c := require.New(t)
	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err)

	err = CheckAltieThemes(tmpDir)
	c.NoError(err)

	// Test case 2: When the directory does not exist
	nonExistingDir := "NonExistingDir"
	err = CheckAltieThemes(nonExistingDir)
	c.Error(err)
}
