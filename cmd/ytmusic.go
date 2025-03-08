package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cookie string

func init() {
	rootCmd.AddCommand(ytCmd)
	ytCmd.Flags().
		StringVarP(&cookie, "cookie-set", "c", "", "Authentication cookie copied from request header")
}

var ytCmd = &cobra.Command{
	Use:   "yt",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(cookie) == 0 {
			return fmt.Errorf("invalid cookie")
		}
		viper.Set("yt-cookie", cookie)
		err := viper.WriteConfig()
		return err
	},
}
