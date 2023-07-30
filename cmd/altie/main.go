package main

import (
	"fmt"
	"os"

	"log"

	"github.com/copydataai/altie/internal/config"
	"github.com/copydataai/altie/internal/themes"
	"github.com/hackebrot/turtle"
	"github.com/pterm/pterm"
)

func ListThemes(homeDir string) {
	dirs, err := themes.ListThemes(homeDir, config.RouteThemes)
	if err != nil {
		pterm.Error.PrintOnError(err)
		return
	}
	inteSelection := pterm.DefaultInteractiveSelect.WithOptions(dirs)

	selectedOption, err := inteSelection.Show()
	if err != nil {
		pterm.Error.PrintOnError(err)
		return
	}

	// TODO: Implement a method to read ThemesDirectory and concatenate them
	path := fmt.Sprintf(config.RouteThemes+"/%s", homeDir, selectedOption)

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

	// TODO: implement a method to use altie.conf
	_, err = config.CheckConfig(homeDir)
	if os.IsNotExist(err) {
		pterm.Printfln("Do you want to create a default altie config in %s/.altie/altie.conf?", homeDir)
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
		_, err = config.CheckConfig(homeDir)
	}

	if err != nil {
		pterm.Error.PrintOnError(err)
		return
	}

	err = themes.CheckAltieThemes(homeDir, config.RouteThemes)
	if os.IsNotExist(err) {
		pterm.Printfln("Do you want to copy all the themes to %s/.altie/themes?", homeDir)
		result, _ := pterm.DefaultInteractiveConfirm.Show()
		if !result {
			pterm.Info.Printfln("You will need to create manual a dir themes with all themes you want in %s/.altie/themes", homeDir)
			return
		}

		configDir := fmt.Sprintf(config.RouteThemes, homeDir)
		themeRepoDir, err := themes.GetRepoDirectory()
		if err != nil {
			pterm.Error.Println("Please Make sure it is in the repository directory")
			pterm.Error.PrintOnError(err)
		}

		err = themes.CreateThemes(configDir, themeRepoDir)
		if err != nil {
			pterm.Error.PrintOnError(err)
			return
		}

		pterm.Info.Printfln("it's created the themes in %s/.altie/themes", homeDir)
	}

	ListThemes(homeDir)
}
