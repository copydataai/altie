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

func ListThemes(altieConfig *config.ConfigThemes, homeDir string) error {
	dirs, err := themes.ListThemes(altieConfig.Config.ThemesDirectory)
	if err != nil {
		return err
	}
	inteSelection := pterm.DefaultInteractiveSelect.WithOptions(dirs)

	selectedOption, err := inteSelection.Show()
	if err != nil {
		return err
	}

	path := fmt.Sprintf(altieConfig.Config.ThemesDirectory+"/%s", selectedOption)
	pterm.Info.Println(path)
	alacrittyConfDir := fmt.Sprintf(config.AlacrittyConfigDir, homeDir)

	backupTheme, err := themes.BackUpTheme(alacrittyConfDir)
	if err != nil {
		return err
	}

	pterm.Info.Printfln("The last theme was saved as %s", backupTheme)

	err = themes.ApplyTheme(path, alacrittyConfDir)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Selected option: %s has been applied successful", pterm.Green(selectedOption))

	return nil
}

func main() {
	pterm.Printfln("Welcome to Altie \nan alternative version of alacritty-themes\nhas been building with Go %s", turtle.Emojis["bear"])
	homeDir, err := config.GetHomeDir()
	if err != nil {
		log.Fatal(err)
		return
	}

	configDir := fmt.Sprintf(config.RouteConfig, homeDir)

	altieConfig, err := config.CheckConfig(configDir)
	if os.IsNotExist(err) {
		pterm.Printfln("Do you want to create a default altie config in %s/.altie/altie.conf?", homeDir)
		result, _ := pterm.DefaultInteractiveConfirm.Show()
		if !result {
			pterm.Info.Printfln("You will need to create manual an altie.config in %s/.altie/altie.conf", homeDir)
			return
		}
		err = config.CreateConfig(homeDir)
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

	err = themes.CheckAltieThemes(altieConfig.Config.ThemesDirectory)
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

	err = ListThemes(altieConfig, homeDir)
	if err != nil {
		pterm.Error.PrintOnError(err)
	}
}
