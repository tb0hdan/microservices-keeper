package input_slack // nolint

import (
	"fmt"
	"os"
	"strings"

	"github.com/tb0hdan/microservices-keeper/repository/logs"

	log2 "log"

	"github.com/nlopes/slack"
)

func RunWebsockets(config *SlackConfiguration) (err error) { // nolint
	api := slack.New(
		config.APIToken,
		slack.OptionDebug(true),
		// FIXME: Switch to internal logs package instead
		slack.OptionLog(log2.New(os.Stdout, "slack-bot: ", log2.Lshortfile|log2.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		logs.Logger.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			logs.Logger.Println("Infos:", ev.Info)
			logs.Logger.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C2147483705"))

		case *slack.MessageEvent:
			logs.Logger.Printf("Message: %v\n", ev)
			msg := strings.Replace(ev.Text, "<!channel>", "@channel", 1)
			// only @channel messages are processed
			if !strings.HasPrefix(msg, "@channel") {
				continue
			}
			// remove @channel
			msg = strings.TrimPrefix(msg, "@channel ")
			reply, err2 := config.MessageHandler(msg)
			if err2 != nil {
				reply = fmt.Sprintf("An error occured while storing message: %+v", err)
				logs.Logger.Printf(reply)

			}
			rtm.SendMessage(rtm.NewOutgoingMessage(reply, ev.Channel))

		case *slack.PresenceChangeEvent:
			logs.Logger.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			logs.Logger.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			logs.Logger.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			logs.Logger.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// log.Printf("Unexpected: %v\n", msg.Data)
		}
	}
	return err
}
