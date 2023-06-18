package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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

	themesDirectory := fmt.Sprintf(config.DirectoryThemes, repoDirectory)
	configDirectory := fmt.Sprintf(config.RouteThemes, homeDir)

	err = CreateThemes(configDirectory, themesDirectory)
	c.NoError(err)

	newDirectory, err := GetRepoDirectory()
	c.NoError(err)

	themesDirectory = fmt.Sprintf(config.DirectoryThemes, newDirectory)

	c.NoError(err)
	err = CreateThemes(configDirectory, themesDirectory)
	c.Error(err)
}
