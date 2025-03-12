package cmd

import (
	"fmt"
	"os"

	"github.com/Samarthbhat52/soundport/cmd/ui/listcommon.go"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(spotifyCmd)
}

type listOptions struct {
	options []string
}

var spotifyCmd = &cobra.Command{
	Use:   "port",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		listOfStuff := []string{
			"Some",
			"Options",
			"Let's flesh it out later",
		}

		p := tea.NewProgram(listcommon.InitialModel(listOfStuff))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas there's been an error: %v\n", err)
			os.Exit(1)
		}
	},
}
