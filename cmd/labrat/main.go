/*
Released under MIT license, copyright 2019 Tyler Ramer
*/

package main

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tylarb/LabRat/pkg/labrat"
)

var rootCmd = &cobra.Command{
	Use:   "labrat",
	Short: "The labrat cli",
	Long:  "The labrat cli: start a container with a tmate session - output is the ssh command used to connect to the session",
}

var (
	tmateConfig      string
	testStringReader = strings.NewReader("labrat cheese")
)

func init() {
	rootCmd.PersistentFlags().StringVar(&tmateConfig, "tmateconf", "", "Use the specified tmate conf")
	rootCmd.AddCommand(labrat.SessionCmd)
	rootCmd.AddCommand(labrat.CheeseCmd)

}

func main() {
	rootCmd.Execute()
}
