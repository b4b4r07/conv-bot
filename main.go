package main

import (
	"log"
	"os"
	"regexp"

	"github.com/nlopes/slack"
)

var pattern *regexp.Regexp = regexp.MustCompile(`^bot\s+test\s+(.*)`)
var (
	asked    bool
	answered bool
)

func main() {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	os.Exit(run(api))
}

func run(api *slack.Client) int {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	var params slack.PostMessageParameters

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Print("Connected!")

			case *slack.MessageEvent:
				pat := pattern.FindStringSubmatch(ev.Text)
				if len(pat) < 2 && !asked {
					break
				}
				if asked {
					switch ev.Text {
					case "":
						break
					case "yes":
						params = getPostMessageParameters("ok", true)
						api.PostMessage(ev.Channel, "", params)
						answered = false
						asked = false
						break
					default:
						params = getPostMessageParameters("canceled", true)
						api.PostMessage(ev.Channel, "", params)
						answered = false
						asked = false
						break
					}
				} else {
					params = getPostMessageParameters("are you ok? y/n", true)
					api.PostMessage(ev.Channel, "", params)
					asked = true
				}

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}
		}
	}
}

func getPostMessageParameters(result string, ok bool) slack.PostMessageParameters {
	color := "danger"
	if ok {
		color = "good"
	}

	params := slack.PostMessageParameters{
		Markdown:  true,
		Username:  "conv-bot",
		IconEmoji: ":trollface:",
	}
	params.Attachments = []slack.Attachment{}
	params.Attachments = append(params.Attachments, slack.Attachment{
		Fallback:   "",
		Title:      "",
		Text:       result,
		MarkdownIn: []string{"title", "text", "fields", "fallback"},
		Color:      color,
	})
	return params
}
