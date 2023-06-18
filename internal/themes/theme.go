package themes

import (
	"os"

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
