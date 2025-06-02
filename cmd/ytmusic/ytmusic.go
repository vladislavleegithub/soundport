package ytmusic

import (
	"fmt"
	"os"
	"strings"

	"github.com/Samarthbhat52/soundport/logger"
	"github.com/Samarthbhat52/soundport/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	glbLogger = logger.GetInstance()
	Cmd       = &cobra.Command{
		Use:   "ytmusic",
		Short: "Handles YouTube Music functionalities.",
	}
)

var ytmusicCookie string

var ytmForm = huh.NewForm(
	huh.NewGroup(
		huh.NewText().
			Title("YT Music cookie").
			Description("Eneter your YT Music cookie").
			Value(&ytmusicCookie).
			CharLimit(2000).
			Validate(huh.ValidateNotEmpty()),
	).Title("YT Music setup"),
)

var ytmusicSetup = &cobra.Command{
	Use:   "setup",
	Short: "Sets up Youtube Music credentials.",
	Long:  "Prompts user to input their Youtube Music authentication cookie extracted from the browser. This ensures that playlist creation has the required credentials.",
	Run: func(cmd *cobra.Command, args []string) {
		var status strings.Builder

		err := ytmForm.Run()
		if err != nil {
			glbLogger.Println(err)
			fmt.Println("Something went wrong")
			os.Exit(1)
		}

		cleanedCookie := strings.ReplaceAll(ytmusicCookie, "\u0026", "")

		viper.Set("yt-cookie", cleanedCookie)
		err = viper.WriteConfig()
		if err != nil {
			glbLogger.Println("Error writing configuration:", err)
			status.WriteString(ui.Red.Render("Error writing configuration\n"))
			os.Exit(1)
		}

		status.WriteString(ui.Green.Render("Setup successful\n"))
		fmt.Println(status.String())
	},
}

func init() {
	Cmd.AddCommand(ytmusicSetup)
}
