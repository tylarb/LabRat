/*
Released under MIT license, copyright 2019 Tyler Ramer
*/

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/tylarb/LabRat/pkg/labrat"
)

var (
	slackBotToken string
	slackBotName  string
	botID         string
)

var (
	sc  *slack.Client
	rtm *slack.RTM
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	configureCmds()
}

func getBotID(botName string, sc *slack.Client) (botID string) {
	users, err := sc.GetUsers()
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		if user.Name == botName {
			log.WithFields(log.Fields{"ID": user.ID, "name": user.Name}).Debug("Found bot:")
			botID = user.ID
		}
	}
	return
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
	sc = slack.New(slackBotToken)
	botID = getBotID(slackBotName, sc)

	rtm = sc.NewRTM()
	go rtm.ManageConnection()
	log.Info("Connected to slack")
	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			log.WithFields(log.Fields{"Channel": ev.Channel, "message": ev.Text}).Debug("message event:")
			if ev.Text == "" {
				continue
			}
			out, err := parse(ev.Text)
			if err != nil {
				log.WithField("ERROR", err).Error("Parsing failed")
				rtm.SendMessage(rtm.NewOutgoingMessage("Unable to create session", ev.Channel))
				continue
			}
			switch {
			case out == "": //continue
			case out == labrat.Cheese:
				printCheese(ev, out)
			default:
				rtm.SendMessage(rtm.NewOutgoingMessage(rawtext(out), ev.Channel))
			}
		case *slack.RTMError:
			log.WithField("ERROR", ev.Error()).Error("RTM Error")

		case *slack.InvalidAuthEvent:
			log.Error("Invalid credentials")
			return

		default:
		}

	}
}

func parse(message string) (string, error) {
	atBot := usrFormat(botID)
	words := strings.Fields(message)
	if words[0] == atBot {
		return execute(words[1:])
	}
	return "", nil

}

func execute(args []string) (string, error) {
	output := new(bytes.Buffer)
	// confirm we clean the buffer for the next command execution
	defer output.Reset()
	labrat.SetOut(output, output)

	if err := labrat.Execute(args); err != nil {
		return "", err
	}
	return output.String(), nil
}

func usrFormat(u string) string {
	return fmt.Sprintf("<@%s>", u)
}

func printCheese(ev *slack.MessageEvent, cheese string) {
	channel, timestamp, _ := sc.PostMessage(ev.Channel, slack.MsgOptionAsUser(true), slack.MsgOptionText("beep beep whirrrree", false))
	linebyline := strings.Split(cheese, "\n")
	completedImage := []string{}
	for _, line := range linebyline {
		time.Sleep(200 * time.Millisecond)
		completedImage = append(completedImage, line)
		image := "```" + strings.Join(completedImage, "\n") + "```"
		sc.UpdateMessage(channel, timestamp, slack.MsgOptionText(image, false))
	}

}

func rawtext(message string) string {
	return "```\n" + message + "\n```"
}
