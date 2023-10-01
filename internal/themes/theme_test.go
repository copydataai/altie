package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestListThemes(t *testing.T) {
	c := require.New(t)

	emptyDir := ""
	expectedEmptyDir := []string{}

	resultEmptyDir, errEmptyDir := ListThemes(emptyDir)
	c.NoError(errEmptyDir)
	c.Equal(expectedEmptyDir, resultEmptyDir)

	// Test case 2: Directory with themes
	dirWithThemes, err := GetRepoDirectory()
	c.NoError(err)

	dirWithThemes = filepath.Join(dirWithThemes, "themes")

	expectedDirWithThemes := "3024.dark.yml"
	resultDirWithThemes, errDirWithThemes := ListThemes(dirWithThemes)
	c.NoError(errDirWithThemes)
	c.Contains(resultDirWithThemes, expectedDirWithThemes)
}

func TestGetRepoDirectory(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err)

	defer os.RemoveAll(tmpDir)

	dir, err := GetRepoDirectory()
	c.NoError(err)
	c.NotEmpty(dir)

	subDir := filepath.Join(tmpDir, "subtest")

	err = os.MkdirAll(subDir, os.ModePerm)
	c.NoError(err)

	err = os.Chdir(subDir)
	c.NoError(err)

	dir, err = GetRepoDirectory()
	c.Error(err)
	c.EqualError(err, ErrNotOnRepoDir.Error())
	c.Empty(dir)

	err = os.RemoveAll(subDir)
	c.NoError(err)

	dir, err = GetRepoDirectory()
	c.Error(err)
	c.EqualError(err, fmt.Errorf("getwd: no such file or directory").Error())
	c.Empty(dir)
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
