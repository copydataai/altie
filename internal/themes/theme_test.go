package themes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Mock for the GithubDirectories  interface.
type MockGithubDirectories struct {
	ListDirectoriesFunc func(url string) ([]themeFile, error)
}

func (mdl *MockGithubDirectories) ListDirectories(url string) ([]themeFile, error) {
	return mdl.ListDirectoriesFunc(url)
}

type MockAltieGithub struct {
	MockDownload func(url string) ([]byte, error)
}

func (m *MockAltieGithub) Download(url string) ([]byte, error) {
	return m.MockDownload(url)
}

type MockAltieTheme struct {
	MockCreateFile func(name string, content []byte, directory string) error
}

func (m *MockAltieTheme) CreateFile(name string, content []byte, directory string) error {
	return m.MockCreateFile(name, content, directory)
}

func TestListThemesOnline(t *testing.T) {
	c := require.New(t)
	mockListDirFunc := func(url string) ([]themeFile, error) {
		return []themeFile{{name: "theme1", url: "http://example.com/theme1"}}, nil
	}
	mockDownloadFunc := func(url string) ([]byte, error) {
		return []byte("theme content"), nil
	}
	mockCreateFileFunc := func(name string, content []byte, directory string) error {
		return nil
	}

	err := ListThemesOnline("/tmp", &MockGithubDirectories{mockListDirFunc}, &MockAltieGithub{mockDownloadFunc}, &MockAltieTheme{mockCreateFileFunc})
	c.NoError(err)

	mockListDirFunc = func(url string) ([]themeFile, error) {
		return nil, errors.New("failed to list directories")
	}
	err = ListThemesOnline("/tmp", &MockGithubDirectories{mockListDirFunc}, &MockAltieGithub{mockDownloadFunc}, &MockAltieTheme{mockCreateFileFunc})
	c.Error(err)
	c.EqualError(err, "failed to list directories")

	mockListDirFunc = func(url string) ([]themeFile, error) {
		return []themeFile{{name: "theme1", url: "http://example.com/theme1"}}, nil
	}
	mockDownloadFunc = func(url string) ([]byte, error) {
		return nil, errors.New("failed to download theme")
	}

	err = ListThemesOnline("/tmp", &MockGithubDirectories{mockListDirFunc}, &MockAltieGithub{mockDownloadFunc}, &MockAltieTheme{mockCreateFileFunc})
	c.Error(err)
	c.EqualError(err, "failed to download theme")

	mockListDirFunc = func(url string) ([]themeFile, error) {
		return []themeFile{{name: "theme1", url: "http://example.com/theme1"}}, nil
	}
	mockDownloadFunc = func(url string) ([]byte, error) {
		return []byte("theme content"), nil
	}
	mockCreateFileFunc = func(name string, content []byte, directory string) error {
		return errors.New("failed to create file")
	}
	err = ListThemesOnline("/tmp", &MockGithubDirectories{mockListDirFunc}, &MockAltieGithub{mockDownloadFunc}, &MockAltieTheme{mockCreateFileFunc})
	c.Error(err)
	c.EqualError(err, "failed to create file")
}

func TestDownloadInsertFiles(t *testing.T) {
	c := require.New(t)

	themes := []themeFile{
		{name: "success.txt", url: "http://example.com/success"},
	}
	mockDownload := func(url string) ([]byte, error) {
		return []byte("mock data"), nil
	}
	mockCreateFile := func(name string, content []byte, directory string) error {
		return nil
	}

	err := downloadInsertFiles(themes, "/tmp", &MockAltieGithub{mockDownload}, &MockAltieTheme{mockCreateFile})
	c.NoError(err)

	themes = []themeFile{
		{name: "fail_download.txt", url: "http://example.com/fail"},
	}
	mockDownload = func(url string) ([]byte, error) {
		return nil, errors.New("download failed")
	}
	mockCreateFile = func(name string, content []byte, directory string) error {
		return nil
	}

	err = downloadInsertFiles(themes, "/tmp", &MockAltieGithub{mockDownload}, &MockAltieTheme{mockCreateFile})
	c.Error(err)
	c.EqualError(err, "download failed")

	themes = []themeFile{
		{name: "fail_create.txt", url: "http://example.com/success"},
	}
	mockDownload = func(url string) ([]byte, error) {
		return []byte("mock data"), nil
	}
	mockCreateFile = func(name string, content []byte, directory string) error {
		return errors.New("file creation failed")
	}
	err = downloadInsertFiles(themes, "/tmp", &MockAltieGithub{mockDownload}, &MockAltieTheme{mockCreateFile})
	c.Error(err)
	c.EqualError(err, "file creation failed")
}

func TestCreateFile(t *testing.T) {
	c := require.New(t)

	tmpDir, err := os.MkdirTemp("", "test")
	c.NoError(err)

	content := []byte("hello text")

	altieTheme := AltieTheme{}

	err = altieTheme.CreateFile("test.txt", content, tmpDir)
	c.NoError(err)

	fContent, err := os.ReadFile(tmpDir + "/test.txt")
	c.NoError(err)
	c.Equal(fContent, content)

	err = os.RemoveAll(tmpDir)
	c.NoError(err)

	err = altieTheme.CreateFile("test.txt", content, tmpDir)
	c.Error(err)
	c.True(os.IsNotExist(err))
}

func TestDownload(t *testing.T) {
	c := require.New(t)
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "successful response")
	}))
	defer successServer.Close()

	// Test server simulating a non-200 response
	non200Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer non200Server.Close()

	github := AltieGithub{}

	body, err := github.Download(successServer.URL)
	c.NoError(err)
	c.Equal(body, []byte("successful response"))

	body, err = github.Download("http://127.0.0.1:0")
	c.Error(err)
	c.Empty(body)

	body, err = github.Download(non200Server.URL)
	c.Error(err)
	c.EqualError(err, ErrCouldNotDownload.Error())
	c.Empty(body)
}

func TestListDirectories(t *testing.T) {
	c := require.New(t)
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]map[string]any{
			{"name": "theme1", "download_url": "http://example.com/theme1"},
			{"name": "theme2", "download_url": "http://example.com/theme2"},
		})
	}))
	defer successServer.Close()

	// Mock server for non-200 response
	non200Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer non200Server.Close()

	// Mock server for invalid JSON response
	invalidJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid JSON"))
	}))
	defer invalidJSONServer.Close()

	// Mock server for missing fields in JSON response
	missingFieldsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]map[string]any{
			{"name": "theme1"},
			{"download_url": "http://example.com/theme2"},
		})
	}))
	defer missingFieldsServer.Close()

	lister := AltieLister{}

	themes, err := lister.ListDirectories(successServer.URL)
	c.NoError(err)
	c.Len(themes, 2)
	c.Equal(themes, []themeFile{{"theme1", "http://example.com/theme1"}, {"theme2", "http://example.com/theme2"}})

	themes, err = lister.ListDirectories(non200Server.URL)
	c.Error(err)
	c.EqualError(err, ErrNotFoundFilesGitHub.Error())

	themes, err = lister.ListDirectories("http://127.0.0.1:0")
	c.Error(err)
	c.EqualError(err, ErrNotFoundFilesGitHub.Error())

	// Error not specifically handled in function
	themes, err = lister.ListDirectories(invalidJSONServer.URL)
	c.Error(err)

	// Technically not an error scenario as per current implementation
	// themes, err = lister.ListDirectories(missingFieldsServer.URL)
	// c.Error(err)
}

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

	// that
	dirWithThemes = filepath.Dir(dirWithThemes)
	dirWithThemes = filepath.Dir(dirWithThemes)

	dirWithThemes = filepath.Join(dirWithThemes, "themes")

	expectedDirWithThemes := "3024.dark.toml"
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

	// create a file with name alacritty.toml inside tmpDir
	alacrittyConfDir := fmt.Sprintf("%s/alacritty.toml", tmpDir)
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
