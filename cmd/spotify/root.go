package spotify

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(setupCmd)
	Cmd.AddCommand(loginCmd)
}

var Cmd = &cobra.Command{
	Use:   "spotify",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
}
