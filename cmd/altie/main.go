package main

import (
	"fmt"
	"os"

	"github.com/copydataai/altie/internal/config"
	"github.com/copydataai/altie/internal/themes"
	"github.com/hackebrot/turtle"
	cp "github.com/otiai10/copy"
	"github.com/pterm/pterm"
)

func ListThemes(altieConfig *config.ConfigThemes, homeDir string) error {
	dirs, err := themes.ListThemes(altieConfig.Config.ThemesDirectory)
	if err != nil {
		return err
	}

	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(dirs).Show()
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

	err = cp.Copy(path, alacrittyConfDir)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Selected option: %s has been applied successful", pterm.Green(selectedOption))

	return nil
}

func CreateConfig() error {
	homeDir, err := config.GetHomeDir()
	if err != nil {
		return err
	}

	configDir := fmt.Sprintf(config.RouteConfig, homeDir)

	altieConfig, err := config.CheckConfig(configDir)
	if os.IsNotExist(err) {
		pterm.Printfln("Do you want to create a default altie config in %s/.altie/altie.conf?", homeDir)
		result, _ := pterm.DefaultInteractiveConfirm.Show()
		if !result {
			pterm.Info.Printfln("You will need to create manual an altie.config in %s/.altie/altie.conf", homeDir)
			// Create a new error when he press CTRL+C
			return nil
		}
		err = config.CreateConfig(homeDir)
		if err != nil {
			return err
		}

		pterm.Info.Printfln("it's created the altie.conf in %s/.altie/altie.conf", homeDir)
		altieConfig, err = config.CheckConfig(configDir)
	}

	if err != nil {
		return err
	}

	err = themes.CheckAltieThemes(altieConfig.Config.ThemesDirectory)
	if os.IsNotExist(err) {
		pterm.Printfln("Do you want to copy all the themes to %s/.altie/themes?", homeDir)
		result, _ := pterm.DefaultInteractiveConfirm.Show()
		if !result {
			pterm.Info.Printfln("You will need to create manual a dir themes with all themes you want in %s/.altie/themes", homeDir)
			return nil
		}

		// Create themes directory
		err := os.MkdirAll(altieConfig.Config.ThemesDirectory, os.ModePerm)
		if err != nil {
			return err
		}

		err = themes.ListThemesOnline(altieConfig.Config.ThemesDirectory, &themes.AltieLister{}, &themes.AltieGithub{}, &themes.AltieTheme{})
		if err != nil {
			return err
		}

		pterm.Info.Printfln("it's created the themes in %s/.altie/themes", homeDir)
	}

	err = ListThemes(altieConfig, homeDir)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	pterm.Printfln("Welcome to Altie \nan alternative version of alacritty-themes\nhas been building with Go %s", turtle.Emojis["bear"])

	err := CreateConfig()
	if err != nil {
		pterm.Error.PrintOnError(err)
	}
}
