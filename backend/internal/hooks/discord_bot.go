package hooks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/meilisearch/meilisearch-go"
	"github.com/pocketbase/pocketbase"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
)

type Vod struct {
	UUID      string     `json:"uuid"`
	Title     string     `json:"title"`
	Date      time.Time  `json:"date"`
	Filename  string     `json:"filename"`
	Viewcount int        `json:"viewcount"`
	Clips     []struct{} `json:"clips"`
}

type SearchResponse struct {
	Error  bool  `json:"error"`
	Result []Vod `json:"result"`
}

type MeiliVod struct {
	UUID      string  `json:"uuid"`
	Title     string  `json:"title"`
	Date      int64   `json:"date"`
	Filename  string  `json:"filename"`
	Viewcount int     `json:"viewcount"`
	Start     float32 `json:"start"`
	Formatted struct {
		Text string `json:"text"`
	} `json:"_formatted"`
}

type MeiliTranscriptResponse struct {
	Hits             []MeiliVod `json:"hits"`
	TotalHits        int        `json:"totalHits"`
	ProcessingTimeMs int        `json:"processingTimeMs"`
}

type UUIDResponse struct {
	Error  bool `json:"error"`
	Result Vod  `json:"result"`
}

type StatsResponse struct {
	CountVods      int     `json:"count_vods"`
	CountClips     int     `json:"count_clips"`
	CountH         float64 `json:"count_hours"`
	CountSizeBytes int     `json:"count_size"`
	TopChatter     []struct {
		Name     string `json:"name"`
		MsgCount int    `json:"msg_count"`
	} `json:"chatters"`
}

var (
	pb             *pocketbase.PocketBase
	token          = os.Getenv("DISCORD_BOT_TOKEN")
	frontendUrl    = os.Getenv("PUBLIC_FRONTEND_URL")
	apiUrl         = os.Getenv("PUBLIC_API_URL")
	meiliUrl       = os.Getenv("PUBLIC_MEILI_URL")
	meiliSearchKey = os.Getenv("PUBLIC_MEILI_SEARCH_KEY")

	// registered commands
	commands = []*discordgo.ApplicationCommand{{
		Name:        "neuestes",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Neuestes Vod anzeigen",
	}, {
		Name:        "suche",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Archiv nach Transcript durchsuchen",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "text",
				Description: "Transcript",
				Type:        discordgo.ApplicationCommandOptionString,
			},
		},
	}, {
		Name:        "stats",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Archiv Statistiken",
	}}

	// handlers for registered commands
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"neuestes": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			logger.Info.Println("Received \"/neuestes\"")
			contentMsg := ":timer: Lade neuestes Vod"
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: contentMsg,
				},
			})
			if err != nil {
				logger.Error.Println(err)
				return
			}
			vods, err := pb.FindRecordsByFilter("vod", "id != ''", "-date", 1, 0)
			if err != nil || len(vods) == 0 {
				contentMsg = ":x: Keine Ergebnisse"
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &contentMsg,
				})
				if err != nil {
					logger.Error.Println(err)
				}
				return
			}
			clips := len(vods[0].ExpandedAll("clip_via_vod"))
			embeds := []*discordgo.MessageEmbed{{
				Title: vods[0].GetString("title"),
				URL:   fmt.Sprintf("%s/vods/%s", frontendUrl, vods[0].Id),
				Image: &discordgo.MessageEmbedImage{
					URL: fmt.Sprintf("%s/vods/%s-lg.webp", apiUrl, vods[0].GetString("filename")),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Datum",
						Value:  fmt.Sprintf("<t:%d>", vods[0].GetDateTime("date").Time().Unix()),
						Inline: true,
					},
					{
						Name:   "Views",
						Value:  vods[0].GetString("viewcount"),
						Inline: true,
					},
					{
						Name:   "Clips",
						Value:  strconv.Itoa(clips),
						Inline: true,
					},
				},
			}}
			contentMsg = ""
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &contentMsg,
				Embeds:  &embeds,
			})
			if err != nil {
				logger.Error.Println(err)
			}
		},
		"suche": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			data := i.ApplicationCommandData()
			contentMsg := ":warning: Suche darf nicht leer sein"
			if len(data.Options) == 0 || data.Options[0].StringValue() == "" {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: contentMsg,
					},
				})
				if err != nil {
					logger.Error.Println(err)
				}
				return
			}
			searchStr := data.Options[0].StringValue()

			logger.Info.Println("Received \"/suche " + searchStr + "\"")

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: ":timer: **Suche nach:** " + searchStr,
				},
			})
			if err != nil {
				logger.Error.Println(err)
				return
			}
			client := meilisearch.New(meiliUrl, meilisearch.WithAPIKey(meiliSearchKey))
			searchRes, err := client.Index("transcripts").Search(searchStr, &meilisearch.SearchRequest{
				HitsPerPage:           5,
				AttributesToHighlight: []string{"text"},
				HighlightPreTag:       "***",
				HighlightPostTag:      "***",
			})
			if err != nil {
				logger.Error.Println(err)
				contentMsg := ":x: Fehler beim Api Request"
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &contentMsg,
				})
				if err != nil {
					logger.Error.Println(err)
				}
				return
			}
			jsonResp, err := json.Marshal(searchRes)
			if err != nil {
				logger.Error.Println(err)
				return
			}
			var transcriptResponse MeiliTranscriptResponse
			err = json.Unmarshal(jsonResp, &transcriptResponse)
			if err != nil {
				logger.Error.Println(err)
				return
			}
			var content string
			if transcriptResponse.TotalHits == 1 {
				content = fmt.Sprintf("# %d Vod gefunden", transcriptResponse.TotalHits)
			} else {
				content = fmt.Sprintf("# %d Vods gefunden", transcriptResponse.TotalHits)
			}
			content += fmt.Sprintf(" (:rocket: %d ms)\n", transcriptResponse.ProcessingTimeMs)
			emotes := []string{
				":one:",
				":two:",
				":three:",
				":four:",
				":five:",
				":six:",
				":seven:",
				":eight:",
				":nine:",
				":keycap_ten:",
			}
			results := transcriptResponse.Hits
			for i, vod := range results {
				content += fmt.Sprintf("\n## %s [**%s**](<%s/vods/%s?t=%.2f>)", emotes[i], vod.Title, frontendUrl, vod.UUID, vod.Start)
				content += fmt.Sprintf("\n:book: %s", vod.Formatted.Text)
				content += fmt.Sprintf("\n:calendar: <t:%d>\n", vod.Date)
			}
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
			if err != nil {
				logger.Error.Println(err)
			}
		},
		"stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			logger.Info.Printf("Received \"/stats\"")
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: ":timer: Lade Statistiken",
				},
			})
			if err != nil {
				logger.Error.Println(err)
				return
			}
			var response StatsResponse
			if err := Stats(&response); err != nil {
				contentMsg := ":x: Fehler beim Api Request"
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &contentMsg,
				})
				if err != nil {
					logger.Error.Println(err)
				}
				return
			}
			content := fmt.Sprintf(`
# :chart_with_upwards_trend: Allgemein
%d Vods
%d Clips
%.2f Stunden gestreamt
%.2fTiB Archivgröße

# :keyboard: Top Chatter
:one: **%s**, %d Nachrichten
:two: **%s**, %d Nachrichten
:three: **%s**, %d Nachrichten
:four: **%s**, %d Nachrichten
:five: **%s**, %d Nachrichten
:six: **%s**, %d Nachrichten
:seven: **%s**, %d Nachrichten
:eight: **%s**, %d Nachrichten`,
				response.CountVods,
				response.CountClips,
				response.CountH,
				float64(response.CountSizeBytes)/1024/1024/1024/1024,
				response.TopChatter[0].Name,
				response.TopChatter[0].MsgCount,
				response.TopChatter[1].Name,
				response.TopChatter[1].MsgCount,
				response.TopChatter[2].Name,
				response.TopChatter[2].MsgCount,
				response.TopChatter[3].Name,
				response.TopChatter[3].MsgCount,
				response.TopChatter[4].Name,
				response.TopChatter[4].MsgCount,
				response.TopChatter[5].Name,
				response.TopChatter[5].MsgCount,
				response.TopChatter[6].Name,
				response.TopChatter[6].MsgCount,
				response.TopChatter[7].Name,
				response.TopChatter[7].MsgCount,
			)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
			if err != nil {
				logger.Error.Println(err)
			}
		},
	}
)

func RunDiscordBot(app *pocketbase.PocketBase) {
	if token == "" {
		logger.Info.Println("Discord token not set. Not starting bot.")
		return
	}

	pb = app

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error.Println(err)
		return
	}

	// debug login info
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Info.Printf("Logged in as: %s#%s", s.State.User.Username, s.State.User.Discriminator)
	})

	// print join server
	s.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		logger.Info.Printf("Joined server: \"%s\"", g.Name)
	})

	// print leave server
	s.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		logger.Warning.Printf("Left server: \"%s\"", g.BeforeDelete.Name)
	})

	// add handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	s.Client.Timeout = time.Duration(time.Duration.Seconds(10))

	// run bot
	if err := s.Open(); err != nil {
		logger.Error.Println(err)
		return
	}

	// register commands
	createdCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", commands)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	for _, cmd := range createdCommands {
		logger.Info.Printf("Registered command: \"%s\"", cmd.Name)
	}

	logger.Info.Printf("Join link: https://discord.com/oauth2/authorize?client_id=%s&scope=applications.commands%%20bot", s.State.User.ID)

	// wait for kill
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	defer s.Close()
	logger.Info.Println("Bot stopped")
}

func Stats(response *StatsResponse) error {
	res, err := http.Get(fmt.Sprintf("%s/stats", apiUrl))
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err := fmt.Errorf("status code was %d, expected 200", res.StatusCode)
		logger.Error.Println(err)
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}
