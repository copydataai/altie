package themes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/copydataai/altie/internal/config"
)

const (
	githubContentDirectory = "https://api.github.com/repos/copydataai/altie/contents/themes"
)

var (
	ErrNotOnRepoDir        = errors.New("you are not on the repo directory")
	ErrNotFoundFilesGitHub = errors.New(fmt.Sprintf("Failed fetching %s ", githubContentDirectory))
	ErrCouldNotDownload    = errors.New("I could download that theme")
)

type GithubDownloader interface {
	Download(url string) ([]byte, error)
}

type GithubDirectories interface {
	ListDirectories(url string) ([]themeFile, error)
}

type ThemeCreator interface {
	CreateFile(name string, content []byte, directory string) error
}

type AltieTheme struct{}
type AltieGithub struct{}
type AltieLister struct{}

type themeFile struct {
	name string
	url  string
}

func ListThemesOnline(themesDirectory string, lister GithubDirectories, downloader GithubDownloader, creator ThemeCreator) error {
	dirNames, err := lister.ListDirectories(githubContentDirectory)
	if err != nil {
		return err
	}

	err = downloadInsertFiles(dirNames, themesDirectory, downloader, creator)
	if err != nil {
		return err
	}

	return nil
}

func downloadInsertFiles(themes []themeFile, themesDirectory string, github GithubDownloader, themeCreator ThemeCreator) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(themes))

	for _, file := range themes {
		wg.Add(1)
		go func(file themeFile) {
			defer wg.Done()
			output, err := github.Download(file.url)
			if err != nil {
				errChan <- err
			}
			err = themeCreator.CreateFile(file.name, output, themesDirectory)
			if err != nil {
				errChan <- err
			}
		}(file)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (at AltieTheme) CreateFile(name string, content []byte, themesDirectory string) error {
	path := fmt.Sprintf(themesDirectory+"/%s", name)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(content); err != nil {
		return fmt.Errorf("writing to file failed %s: %w", path, err)
	}

	return nil
}

func (ag AltieGithub) Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrCouldNotDownload
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (al AltieLister) ListDirectories(url string) ([]themeFile, error) {
	themesLinks := make([]themeFile, 0)
	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrNotFoundFilesGitHub
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrNotFoundFilesGitHub
	}

	var themesGithub []map[string]any

	err = json.NewDecoder(resp.Body).Decode(&themesGithub)
	if err != nil {
		return nil, err
	}

	for _, item := range themesGithub {
		nameItem, ok := item["name"]
		if !ok {
			continue
		}

		downloadItem, ok := item["download_url"]
		if !ok {
			continue
		}
		themesLinks = append(themesLinks, themeFile{
			name: string(nameItem.(string)),
			url:  string(downloadItem.(string)),
		})
	}

	return themesLinks, nil
}

func GetRepoDirectory() (string, error) {
	dirPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dirPath, nil
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

func ApplyFontTheme(pathConfig string, themeConfig *config.ThemeConfig) error {
	font := themeConfig.Font
	fontSize := themeConfig.FontSize

	alacrittyConfig, err := CheckAlacrittyConfig(pathConfig)
	if err != nil {
		return err
	}

	// TODO: change the map model to struct model
	alacrittyConfig["font"] = map[string]any{
		"normal": map[string]string{
			"family": font,
		},
		"bold": map[string]string{
			"family": font,
		},
		"italic": map[string]string{
			"family": font,
		},
		"bold_italic": map[string]string{
			"family": font,
		},
		"size": fontSize,
	}

	alacrittyFile, err := os.OpenFile(pathConfig, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	err = toml.NewEncoder(alacrittyFile).Encode(alacrittyConfig)
	if err != nil {
		return fmt.Errorf("failed to encode TOML config: %w", err)
	}

	return nil
}

func CheckAlacrittyConfig(pathConfig string) (map[string]any, error) {
	config := make(map[string]any)
	if _, err := toml.DecodeFile(pathConfig, &config); err != nil {
		return config, err
	}

	return config, nil
}
