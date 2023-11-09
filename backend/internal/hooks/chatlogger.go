package hooks

import (
	"html"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

type Chatlogger struct {
	Client     *twitch.Client
	App        *pocketbase.PocketBase
	Collection *models.Collection
}

// Creates a new chatlogger instance
func NewChatlogger(app *pocketbase.PocketBase) (*Chatlogger, error) {
	var err error
	cl := &Chatlogger{}
	cl.App = app
	cl.Client = twitch.NewAnonymousClient()
	cl.Collection, err = app.Dao().FindCollectionByNameOrId("chatmessage")
	if err != nil {
		logger.Error.Println(err)
		return cl, err
	}

	return cl, nil
}

// Run the chatlogger
func (cl *Chatlogger) Run(broadcaster string) {
	logger.Debug.Println("[hooks] running chatlogger")

	cl.Client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		record := models.NewRecord(cl.Collection)
		record.Set("date", message.Time)
		record.Set("user_id", message.User.ID)
		record.Set("user_display_name", message.User.DisplayName)
		record.Set("user_name", message.User.Name)
		record.Set("message", html.EscapeString(message.Message))
		record.Set("tags", message.Tags)
		if err := cl.App.Dao().SaveRecord(record); err != nil {
			logger.Error.Println(err)
		}
	})

	cl.Client.Join(broadcaster)

	for {
		err := cl.Client.Connect()
		if err != nil {
			logger.Error.Println(err)
			time.Sleep(1 * time.Second) // avoid rate limit
		}
	}
}
