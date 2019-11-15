/*
Released under MIT license, copyright 2019 Tyler Ramer
*/

package labrat

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// SessionTimeout is the session timeout in hours
	SessionTimeout time.Duration
	// WaitReadyTimeout is the time to wait for tmate to form a valid connection in seconds
	WaitReadyTimeout = time.Second * 3
)

var rootCmd = &cobra.Command{
	Use:   "labrat",
	Short: "Get the rat to manage the lab",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// SessionCmd gets a new session
var SessionCmd = &cobra.Command{
	Use:     "session [flags]",
	Short:   "Start a new lab session",
	Long:    "The session command allows you to start a new timed ssh session to the lab.\nDefaults to 1 hour session, but can set up to an 8hr session",
	Example: "labrat session -t 4",
	RunE: func(cmd *cobra.Command, args []string) error {
		return CreateSession()
	},
}

// CheeseCmd gives the rat some cheese
var CheeseCmd = &cobra.Command{
	Use:   "cheese",
	Short: "Give the rat something to eat",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), Cheese)
	},
}

func init() {
	rootCmd.AddCommand(SessionCmd)
	rootCmd.AddCommand(CheeseCmd)
	SessionTimeout = time.Hour * time.Duration(*SessionCmd.Flags().Int64P("timeout", "t", 1, "Session timeout, default 1 hr, max of 8 hr"))
}

// SetOut set's labrat's out and err.
// Use it to utilize the cobra commands without requiring
// stdout and stderr from a terminal
func SetOut(outWriter, errWriter io.Writer) {
	commands := []*cobra.Command{
		rootCmd,
		SessionCmd,
		CheeseCmd,
	}
	for _, cmd := range commands {
		cmd.SetOut(outWriter)
		cmd.SetErr(errWriter)
	}
}

func CreateSession() error {
	// Run podman tmate-client
	podmanRun := []string{"podman", "run", "-d", "tmate-client"}
	cmd := exec.Command(podmanRun[0], podmanRun[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	containerID := strings.Trim(string(out[:]), "\n")

	// Go ahead and start the timer now that we've got a container
	go KillSession(containerID, SessionTimeout)

	log.WithField("ctrID", containerID).Info("Container running tmate running")
	// TODO. fix this to a reasonable timeout, but check to see if the session exists
	// with tmate -S /tmp/tmate.sock wait tmate-ready &&
	time.Sleep(time.Second * time.Duration(WaitReadyTimeout))
	getTmateSSH := strings.Fields("tmate -S /tmp/tmate.sock display -p '#{tmate_ssh}'")
	podmanExec := append([]string{"podman", "exec", containerID}, getTmateSSH...)
	cmd = exec.Command(podmanExec[0], podmanExec[1:]...)
	out, err = cmd.CombinedOutput()
	sshSession := string(out[:])
	if err != nil {
		log.WithField("output", sshSession).Error("error on getting tmate session")
		return err
	}
	log.WithField("Session", sshSession).Info("tmate ssh session received")

	fmt.Fprintf(rootCmd.OutOrStdout(), sshSession)
	return nil
}

func KillSession(ctrID string, timeout time.Duration) {
	podmanStop := []string{"podman", "stop", "-t0", ctrID}
	time.Sleep(timeout)
	cmd := exec.Command(podmanStop[0], podmanStop[1:]...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithField("ctr", out[:]).Info("Reaped container, ending tmate session")
	} else {
		log.WithError(err).WithField("ctr", out[:]).Error("Error reaping container")
	}
}

// Execute replaces the os.Args[] with custom args and executes the commands accordingly.
// No passed args executes the root function (i.e. printing help)
func Execute(args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

// Cheese for the rat
const Cheese = `
           ()()
           (..)
           /\/\
    ___ __c\db/o____    _______
 .-" _ "             "=-"  \\   \
|   ( )     _           o   :.   .
|    "     ( )     ()       ::   :
|_          "          ..   ::   :
  )              ()   (  )  :|   |
|"    ()               ""   :|   |
|   O        o .-.     _    :.   /
\____.--._____(---)___(-)__//___/
`
