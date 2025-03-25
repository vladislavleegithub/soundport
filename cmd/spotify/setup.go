package spotify

import (
	"fmt"
	"os"
	"strings"

	"github.com/Samarthbhat52/soundport/cmd/ui"
	"github.com/Samarthbhat52/soundport/logger"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var glbLogger = logger.GetInstance()

var (
	spfyClientId     string
	spfyClientSecret string
)

var spfyForm = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("Client ID").
			Description("Eneter your spotify client ID").
			Value(&spfyClientId).
			Validate(huh.ValidateNotEmpty()),

		huh.NewInput().
			Title("Client Secret").
			Description("Eneter your spotify client secret").
			Value(&spfyClientSecret).
			Validate(huh.ValidateNotEmpty()),
	).Title("Spotify setup").WithWidth(20),
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up Spotify api credentials.",
	Long:  "Prompts user to input their Spotify developer credentials to ensure soundport can interact with Sotify's API.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var status strings.Builder

		err := spfyForm.Run()
		if err != nil {
			glbLogger.Println(err)
			os.Exit(1)
		}

		viper.Set("spfy-id", spfyClientId)
		viper.Set("spfy-secret", spfyClientSecret)

		err = viper.WriteConfig()
		if err != nil {
			glbLogger.Println("Error writing to config")

			status.WriteString(ui.Red.Render("Setup failed\n"))
			fmt.Println(status.String())

			os.Exit(1)
		}

		status.WriteString(ui.Green.Render("Setup successful\n"))
		fmt.Println(status.String())
	},
}
