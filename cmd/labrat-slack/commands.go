/*
Released under MIT license, copyright 2019 Tyler Ramer
*/

package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "labrat-slack [flags]",
	Short: "Start the labrat slack engine",
	Long:  "Provided a bot username and slack OAuth token, connect to slack and run labrat via slack chat",
	RunE: func(cmd *cobra.Command, args []string) error {
		if slackBotName == "" || slackBotToken == "" {
			return errors.New("Please ensure both name and token are populated")
		}
		return nil
	},
}

var (
	tmateConfig string
)

func configureCmds() {
	rootCmd.PersistentFlags().StringVar(&tmateConfig, "tmateconf", "", "Use the specified tmate conf")
	rootCmd.PersistentFlags().StringVarP(&slackBotName, "name", "n", os.Getenv("SLACK_BOT_NAME"), "Use the specified tmate conf")
	rootCmd.PersistentFlags().StringVarP(&slackBotToken, "token", "t", os.Getenv("SLACK_BOT_TOKEN"), "Use the specified tmate conf")

}
