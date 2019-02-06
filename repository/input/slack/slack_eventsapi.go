package input_slack // nolint

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus" // nolint

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

func RunEvents(config *SlackConfiguration) (err2 error){ // nolint
	// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
	var api = slack.New(config.APIToken)

	http.HandleFunc(config.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			log.Printf("body read failed with: %+v\n", err)
		}
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body),
			slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: config.VerificationToken}))
		if e != nil {
			log.Printf("Verification failed with: %+v\n", e)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			_, err = w.Write([]byte(r.Challenge))
			if err != nil {
				log.Printf("body write failed with: %+v\n", err)
			}
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) { // nolint
			case *slackevents.AppMentionEvent:
				msg, err := config.MessageHandler(ev.Text)
				if err == nil && msg != "" {
					_, _, err = api.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
					if err != nil {
						log.Printf("Could not send message to Slack")
					}

				} else {
					log.Printf("Did not send message: %+v, %+v\n", msg, err)
				}

			}
		}
	})
	addr := ":3000"
	log.Printf("[INFO] Server listening on addr: %s", addr)
	if err2 = http.ListenAndServe(addr, nil); err2 != nil {
		log.Printf("HTTP listener error: %+v\n", err2)
	}
	return err2
}
