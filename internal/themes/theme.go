package themes

import (
	"fmt"
	"os"
	"time"

	cp "github.com/otiai10/copy"
)

func GetRepoDirectory() (string, error) {
	dirPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func CreateThemes(configDirectory, themesDirectory string) error {
	err := cp.Copy(themesDirectory, configDirectory)
	if err != nil {
		return err
	}

	return nil
}

func ApplyTheme(pathTheme, alacrittyConfDir string) error {
	err := cp.Copy(pathTheme, alacrittyConfDir)
	if err != nil {
		return err
	}
	return nil
}

func BackUpTheme(alacrittyConfDir string) (string, error) {
	year, month, day := time.Now().Date()
	backupPath := fmt.Sprintf("%s.%d%d%d.bak", alacrittyConfDir, year, month, day)
	err := os.Rename(alacrittyConfDir, backupPath)
	if err != nil {
		return "", err
	}

	return backupPath, nil
}
