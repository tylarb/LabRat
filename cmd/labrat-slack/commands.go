/*
Released under MIT license, copyright 2019 Tyler Ramer
*/

package main

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
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
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	},
}

var (
	tmateConfig string
)

func configureCmds() {
	rootCmd.PersistentFlags().StringVar(&tmateConfig, "tmateconf", "", "Use the specified tmate conf")
	rootCmd.PersistentFlags().StringVarP(&slackBotName, "name", "n", os.Getenv("SLACK_BOT_NAME"), "bot name")
	rootCmd.PersistentFlags().StringVarP(&slackBotToken, "token", "t", os.Getenv("SLACK_BOT_TOKEN"), "bot slack token")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Set logs to verbose/debug")

}
