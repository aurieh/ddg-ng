package searchplugin

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/aurieh/ddg-ng/commandclient"
	"github.com/aurieh/ddg-ng/duckduckgo"
	"github.com/aurieh/ddg-ng/htmlmeta"
	"github.com/aurieh/ddg-ng/utils"
	"github.com/bwmarrin/discordgo"
	"net/url"
	"strings"
)

const truncateSuffix = "..."

func truncate(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length-len(truncateSuffix)] + truncateSuffix
}

func getAnswerURL(answer *duckduckgo.Answer) string {
	if answer.AbstractURL != "" {
		return answer.AbstractURL
	} else if answer.DefinitionURL != "" {
		return answer.DefinitionURL
	} else if answer.Redirect != "" {
		return answer.Redirect
	} else if answer.AbstractSource != "" {
		return answer.AbstractSource
	}
	return ""
}

func getAnswerTitle(query string, answer *duckduckgo.Answer) string {
	if answer.Answer != "" {
		return answer.Answer
	} else if answer.Heading != "" {
		return answer.Heading
	}
	return query
}

var searchEngines = []string{"bing.com", "google.", "startpage.com", "duckduckgo.com", "qwant.com"}

func isSearchEngine(url string) bool {
	for _, sURL := range searchEngines {
		if strings.Contains(url, sURL) {
			return true
		}
	}
	return false
}

func addAnswerMeta(url string, embed *discordgo.MessageEmbed) error {
	// Don't bully search engines
	if isSearchEngine(url) {
		return nil
	}

	metaparser, err := htmlmeta.New(utils.Client, url)
	if err != nil {
		return err
	}
	if embed.Title == "" {
		embed.Title = truncate(metaparser.GetTitle(), 256)
	}
	if embed.Description == "" {
		embed.Description = truncate(metaparser.GetMeta("description", "content"), 2048)
	}
	if image := metaparser.GetOGPMeta("image"); image != "" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: image,
		}
	}
	if url := metaparser.GetOGPMeta("url"); url != "" {
		embed.URL = url
	}
	return nil
}

func createAnswerEmbed(query string, answer *duckduckgo.Answer, embed *discordgo.MessageEmbed) error {
	// TODO: Clean up this logic
	if url := getAnswerURL(answer); url != "" {
		embed.URL = url
		err := addAnswerMeta(url, embed)
		if err != nil {
			return err
		}
	}

	if answer.Definition != "" && embed.Description == "" {
		embed.Description = truncate(answer.Definition, 2048)
	}

	if answer.Image != "" && embed.Image == nil {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: answer.Image,
		}
	}

	var topics []duckduckgo.Topic
	if len(answer.RelatedTopics) > 0 {
		topics = answer.RelatedTopics
	} else if len(answer.Results) > 0 {
		topics = answer.Results
	}

	if len(topics) > 0 {
		if embed.Description != "" {
			embed.Description += "\n\nSee also:"
		} else {
			embed.Description = "See also:"
		}
		fields := []*discordgo.MessageEmbedField{}
		topic := topics[0]
		if topic.FirstURL != "" && embed.URL == "" {
			embed.URL = topic.FirstURL
			err := addAnswerMeta(topic.FirstURL, embed)
			if err != nil {
				return err
			}
		}
		if topic.Icon.URL != "" && embed.Image == nil {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: topic.Icon.URL,
			}
		}
		for _, topic := range topics[1:] {
			if len(embed.Fields) > 24 {
				break
			}
			field := &discordgo.MessageEmbedField{}
			if topic.FirstURL != "" {
				field.Value = truncate(topic.FirstURL, 1024)
			} else {
				field.Value = "[no url]"
			}
			if topic.Text != "" {
				field.Name = truncate(topic.Text, 256)
				fields = append(fields, field)
			}
		}
		embed.Fields = fields
	}
	if embed.URL == "" {
		embed.URL = "https://duckduckgo.com/c/" + url.QueryEscape(query)
	}
	if embed.Title == "" {
		embed.Title = truncate(getAnswerTitle(query, answer), 256)
	}

	return nil
}

// SendQueryResult sends an answer embed as a response to specified message
func SendQueryResult(s *discordgo.Session, m *discordgo.MessageCreate, query string) {
	channel, err := s.State.Channel(m.ChannelID)
	isNSFW := false
	if err != nil {
		log.WithError(err).WithField("channel", m.ChannelID).Errorln("failed to get channel from state")
		return
	} else if channel != nil {
		isNSFW = channel.NSFW
	}
	answer, err := duckduckgo.Query(utils.Client, "\\"+query, "ddg-ng/0.1", !isNSFW, true)
	if err != nil {
		log.WithError(err).Errorln("failed to fetch data from ddg")
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I've encountered an error while trying to access the DDG API: %s", err.Error())) // nolint: vetshadow
		if err != nil {
			log.WithError(err).Errorln("failed to send error message")
		}
		return
	}
	response := &discordgo.MessageEmbed{}
	if err = createAnswerEmbed(query, answer, response); err != nil {
		log.WithError(err).WithField("answer", answer).Errorln("failed to create an answer embed")
		return
	}
	if _, err = s.ChannelMessageSendEmbed(m.ChannelID, response); err != nil {
		log.WithError(err).WithField("answer", answer).Errorln("failed to send a ddg answer embed")
	}
}

// SearchCommand DDG search
func SearchCommand(ctx *commandclient.Context) {
	SendQueryResult(ctx.Session, ctx.Message, strings.Join(ctx.Args, " "))
}
