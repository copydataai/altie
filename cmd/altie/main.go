package main

import (
	"fmt"
	"os"
	"path/filepath"

	"log"

	"github.com/copydataai/altie/internal/config"
	"github.com/copydataai/altie/internal/themes"
	"github.com/hackebrot/turtle"
	"github.com/pterm/pterm"
)

func defaultMain(homeDir string) {
	dirs := make([]string, 0)
	pathDirs := make(map[string]string, 0)
	themesConfig := fmt.Sprintf(config.RouteThemes, homeDir)
	err := filepath.Walk(themesConfig, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			pterm.Error.PrintOnError(err)
		}
		nameFile := info.Name()
		if nameFile == "themes" {
			return nil
		}

		dirs = append(dirs, nameFile)
		pathDirs[nameFile] = path

		return nil
	})
	if err != nil {
		pterm.Error.PrintOnError(err)
		return
	}

	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(dirs).Show()
	if err != nil {
		pterm.Error.PrintOnError(err)
		return
	}

	path, _ := pathDirs[selectedOption]

	pterm.Info.Println(path)
	alacrittyConfDir := fmt.Sprintf(config.AlacrittyConfigDir, homeDir)
	backupTheme, err := themes.BackUpTheme(alacrittyConfDir)
	if err != nil {
		pterm.Error.PrintOnError(err)
		return
	}
	pterm.Info.Printfln("The last theme was saved as %s", backupTheme)

	err = themes.ApplyTheme(path, alacrittyConfDir)
	if err != nil {
		pterm.Error.PrintOnError(err)
		return
	}
	pterm.Success.Printfln("Selected option: %s has been applied successful", pterm.Green(selectedOption))
}

func createConfig(homeDir string) error {
	err := config.CreateConfig(homeDir)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	pterm.Printfln("Welcome to Altie \nan alternative version of alacritty-themes\nhas been building with Go %s", turtle.Emojis["bear"])
	homeDir, err := config.GetHomeDir()
	if err != nil {
		log.Fatal(err)
		return
	}

	configThemes, err := config.CheckConfig(homeDir)
	if os.IsNotExist(err) {
		pterm.Printfln("Do you want to create a default altie config in %s/.altie/altie.conf", homeDir)
		result, _ := pterm.DefaultInteractiveConfirm.Show()
		if !result {
			pterm.Info.Printfln("You will need to create manual an altie.config in %s/.altie/altie.conf", homeDir)
			return
		}
		err = createConfig(homeDir)
		if err != nil {
			pterm.Error.PrintOnError(err)
			return
		}
		pterm.Info.Printfln("it's created the altie.conf in %s/.altie/altie.conf", homeDir)
		configThemes, err = config.CheckConfig(homeDir)
	}

	if err != nil {
		pterm.Error.PrintOnError(err)
		return
	}

	pterm.Print(configThemes)

	// repoDirectory, err := themes.GetRepoDirectory()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// themesDirectory := fmt.Sprintf(config.DirectoryThemes, repoDirectory)
	// configDirectory := fmt.Sprintf(config.RouteThemes, homeDir)
	// err = themes.CreateThemes(configDirectory, themesDirectory)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

}
