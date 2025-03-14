package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	spfyClientId     string
	spfyClientSecret string
	ytmCookie        string
)

var form = huh.NewForm(
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
	).Title("Spotify setup"),

	huh.NewGroup(
		huh.NewText().
			Title("YT Music cookie").
			Description("Eneter your YT Music cookie").
			CharLimit(1810).
			Validate(huh.ValidateNotEmpty()),
	).Title("YT Music setup"),
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var status strings.Builder

		err := form.Run()
		if err != nil {
			glbLogger.Println(err)
			os.Exit(1)
		}

		viper.Set("spfy-id", spfyClientId)
		viper.Set("spfy-secret", spfyClientSecret)
		viper.Set("yt-cookie", ytmCookie)

		err = viper.WriteConfig()
		if err != nil {
			glbLogger.Println("Error writing to config")

			status.WriteString(red.Render("Setup failed\n"))
			fmt.Println(status.String())

			os.Exit(1)
		}

		status.WriteString(green.Render("Setup successful\n"))
		fmt.Println(status.String())
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
