package themes

import (
	"fmt"
	"os"
	"path/filepath"
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

func ListThemes(dirThemes string) ([]string, error) {
	dirs := make([]string, 0)
	err := filepath.Walk(dirThemes, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		nameFile := info.Name()
		if nameFile == "themes" {
			return nil
		}

		dirs = append(dirs, nameFile)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirs, nil
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

func CheckAltieThemes(dirThemes string) error {
	_, err := os.Stat(dirThemes)
	if err != nil {
		return err
	}

	return nil
}
