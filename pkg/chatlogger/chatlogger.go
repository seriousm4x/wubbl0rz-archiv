package chatlogger

import (
	"html"
	"os"

	"github.com/AgileProggers/archiv-backend-go/pkg/database"
	"github.com/AgileProggers/archiv-backend-go/pkg/logger"
	"github.com/AgileProggers/archiv-backend-go/pkg/models"
	"github.com/gempir/go-twitch-irc/v3"
)

func Run() {
	client := twitch.NewAnonymousClient()
	var msg models.ChatMessage

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		msg.ID = message.ID
		msg.CreatedAt = message.Time
		msg.UserID = message.User.ID
		msg.UserDisplayName = message.User.DisplayName
		msg.UserName = message.User.Name
		msg.Message = html.EscapeString(message.Message)
		msg.Tags = message.Tags
		if result := database.DB.Model(&msg).Create(msg); result.Error != nil {
			logger.Error.Println(result.Error)
			return
		}
	})

	client.Join(os.Getenv("BROADCASTER_NAME"))

	for {
		err := client.Connect()
		if err != nil {
			logger.Error.Println(err)
		}
	}
}
