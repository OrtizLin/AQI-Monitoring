package linebot

import (
	"aqiCrawler/db"
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
		case linebot.EventTypeFollow:
			if err := app.handleFollow(event.ReplyToken, event.Source); err != nil {
				log.Print(err)
			}
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
}

func (app *LineBotStruct) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "查詢":
		str := db.CheckRegistered(source.UserID)
		replymessage := ""
		if str == "" {
			replymessage = "您目前尚未新增任何觀測站。您可以點擊快速選單裡的註冊觀測站來新增你附近的觀測站"
		} else if str == "nouser" {
			replymessage = "您尚未使用 Line Notify，請先點擊下列網址與 Line Notify 連動\nhttps://aqi-push-notify.herokuapp.com/auth?client=" + source.UserID
		} else {
			replymessage = "目前加入的觀測站有 : " + str
		}
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(replymessage),
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
	if db.NewSite(str, source.UserID) {
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("已將觀測站 : ["+str+"] 加入您的資料庫, 若該地區空氣良好時將主動通知您。"),
		).Do(); err != nil {
			return err
		}
	} else {
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("觀測站加入失敗，請先點擊下列網址與 Line Notify 連動\nhttps://aqi-push-notify.herokuapp.com/auth?client="+source.UserID),
		).Do(); err != nil {
			return err
		}
	}
	return nil
}

func (app *LineBotStruct) handleFollow(replyToken string, source *linebot.EventSource) error {
	//GET USER PROFILE
	profile, err := app.bot.GetProfile(source.UserID).Do()
	if err != nil {
		log.Print(err)
	}
	if _, err := app.bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(profile.DisplayName+" 歡迎使用 空汙報報 。\n請按步驟啟用 LINE Notify\n已獲得最新文章通知。\n1. 開啟下方網址\n2. 選擇“群組”第一個「透過1對1聊天」接收 Line Notify 的通知。\n3. 點擊「同意並連動」。\nhttps://aqi-push-notify.herokuapp.com/auth?client="+profile.UserID),
	).Do(); err != nil {
		return err
	}
	return nil
}
