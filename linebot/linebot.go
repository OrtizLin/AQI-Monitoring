package linebot

import (
	"aqiCrawler/distance"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"strconv"
)

type LineBotStruct struct {
	bot         *linebot.Client
	appBaseURL  string
	downloadDir string
}

func NewLineBot(channelSecret, channelToken, appBaseURL string) (*LineBotStruct, error) {
	bots, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return &LineBotStruct{
		bot:         bots,
		appBaseURL:  appBaseURL,
		downloadDir: "testing",
	}, nil
}

func (app *LineBotStruct) Callback(w http.ResponseWriter, r *http.Request) {
	events, err := app.bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := app.handleText(message, event.ReplyToken, event.Source); err != nil {
					log.Print(err)
				}
			case *linebot.LocationMessage:
				if err := app.handleLocation(message, event.ReplyToken, event.Source); err != nil {
					log.Print(err)
				}
			default:
				log.Printf("Unknown message: %v", message)
			}
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
}

func (app *LineBotStruct) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "tonygrr":
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("嫩！"),
		).Do(); err != nil {
			return err
		}
	default:

		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("嫩！"),
		).Do(); err != nil {
			return err
		}
	}
	return nil

}

func (app *LineBotStruct) handleLocation(message *linebot.LocationMessage, replyToken string, source *linebot.EventSource) error {

	lat := strconv.FormatFloat(message.Latitude, 'f', -1, 64)
	long := strconv.FormatFloat(message.Longitude, 'f', -1, 64)
	str := distance.GetSite(lat, long)

	//GET USER PROFILE
	profile, err := app.bot.GetProfile(source.UserID).Do()
	if err != nil {
		log.Print(err)
	}

	if _, err := app.bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage("離你最近的觀測站為 : "+str+"\n"+profile.DisplayName+"\n"+profile.PictureURL+"\n"+profile.UserID),
	).Do(); err != nil {
		return err
	}
	return nil
}
