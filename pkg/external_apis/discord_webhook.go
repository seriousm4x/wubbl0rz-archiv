package external_apis

// helpful webhook guide: https://birdie0.github.io/discord-webhooks-guide/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AgileProggers/archiv-backend-go/pkg/logger"
)

type discordEmbedImage struct {
	URL string `json:"url"`
}

type discordEmbedThumbnail struct {
	URL string `json:"url"`
}

type discordEmbedAuthor struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

type discordEmbed struct {
	Color       uint                  `json:"color"`
	Author      discordEmbedAuthor    `json:"author"`
	Title       string                `json:"title"`
	URL         string                `json:"url"`
	Description string                `json:"description"`
	Image       discordEmbedImage     `json:"image"`
	Thumbnail   discordEmbedThumbnail `json:"thumbnail"`
	Timestamp   time.Time             `json:"timestamp"`
}

type discordWebhook struct {
	Content string         `json:"content"`
	Embeds  []discordEmbed `json:"embeds"`
}

func DiscordSendWebhook(d TwitchStreamResponse) error {
	logger.Debug.Println("[external_apis] send discord webhook")
	var webhookData discordWebhook
	var embed discordEmbed
	var helixUser TwitchHelixUserResponse

	// get user info from twitch
	if err := TwitchGetHelixUser(&helixUser); err != nil {
		return err
	}

	// prepare webhook data
	streamData := d.Data[0]

	imgUrl := strings.Replace(streamData.ThumbnailURL, "{width}", "1920", 1)
	imgUrl = strings.Replace(imgUrl, "{height}", "1080", 1)

	webhookData.Content = fmt.Sprintf("Auf gehts @everyone, %s macht nun Streamstelz", streamData.UserName)
	embed.Color = 15896107 // hex as decimal
	embed.Author.Name = streamData.UserName
	embed.Author.URL = fmt.Sprintf("https://twitch.tv/%s", streamData.UserLogin)
	embed.Author.IconURL = helixUser.Data[0].ProfileImageUrl
	embed.Title = streamData.Title
	embed.URL = fmt.Sprintf("https://twitch.tv/%s", streamData.UserLogin)
	embed.Description = fmt.Sprintf("%s, %d Viewer", streamData.GameName, streamData.ViewerCount)
	embed.Image.URL = imgUrl
	embed.Thumbnail.URL = fmt.Sprintf("https://static-cdn.jtvnw.net/ttv-boxart/%s-144x192.jpg", streamData.GameID)
	embed.Timestamp = streamData.StartedAt
	webhookData.Embeds = append(webhookData.Embeds, embed)

	// send webhook,
	jsonData, err := json.Marshal(webhookData)
	if err != nil {
		return err
	}

	resp, err := http.Post(os.Getenv("DISCORD_WEBHOOK"), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("discord webhook: status code was %d", resp.StatusCode)
	}

	return nil
}
