package commands

import (
	"github.com/aurieh/ddg-ng/htmlmeta"
	"strconv"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aurieh/ddg-ng/commandclient"
	"github.com/aurieh/ddg-ng/duckduckgo"
	"github.com/aurieh/ddg-ng/github"
	"github.com/aurieh/ddg-ng/stats"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 2,
}

func StatsCommand(ctx *commandclient.Context) {
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, stats.GetStatsString())
}

// func createTopicEmbed(topic duckduckgo.Topic) *discordgo.MessageEmbed {
// 	embed := &discordgo.MessageEmbed{}
// 	if topic.Icon.URL != "" {
// 		embed.Image = &discordgo.MessageEmbedImage{
// 			URL: topic.Icon.URL,
// 		}
// 	}
// 	if topic.FirstURL != "" && topic.Text != "" {
// 		embed.URL = topic.FirstURL
// 		embed.Description = topic.Text
// 	} else if topic.FirstURL != "" {
// 		embed.Description = topic.FirstURL
// 	}
// 	return embed
// }

func createAnswerEmbed(answer *duckduckgo.Answer, embed *discordgo.MessageEmbed) {
	// TODO: Clean up this logic
	if answer.AbstractURL != "" {
		embed.URL = answer.AbstractURL
	} else if answer.DefinitionURL != "" {
		embed.URL = answer.DefinitionURL
	} else if answer.Redirect != "" {
		embed.URL = answer.Redirect
	} else if answer.AbstractSource != "" {
		embed.URL = answer.AbstractSource
	}

	if answer.Answer != "" {
		embed.Title = answer.Answer
	} else if answer.Heading != "" {
		embed.Title = answer.Heading
	} else if embed.URL != "" {
		metaparser, err := htmlmeta.New(client, embed.URL)
		if err != nil {
			log.WithError(err).WithField("url", embed.URL).Errorln("couldnt create htmlmeta for url")
		} else {
			title := metaparser.GetTitle()
			if title != "" {
				embed.Title = title
			} else {
				embed.Title = "[no title in head element]"
			}
			description := metaparser.GetMeta("description", "content")
			if description != "" {
				embed.Description = description
			} else {
				embed.Description = "[no description]"
			}
			url := metaparser.GetOGPMeta("url")
			if url != "" && embed.URL != url {
				embed.URL = url
			}
			image := metaparser.GetOGPMeta("image")
			if image != "" {
				embed.Image = &discordgo.MessageEmbedImage{
					URL: image,
				}
			}
		}
	} else {
		embed.Title = "[no title]"
	}

	if answer.Definition != "" {
		embed.Description = answer.Definition
	}

	if answer.Image != "" {
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

	if topics != nil {
		if embed.Description != "" {
			embed.Description += "\n\nSee also:"
		} else {
			embed.Description = "See also:"
		}
		fields := []*discordgo.MessageEmbedField{}
		topic := topics[0]
		if topic.FirstURL != "" && embed.URL == "" {
			embed.URL = topic.FirstURL
		}
		if topic.Icon.URL != "" && embed.Image == nil {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: topic.Icon.URL,
			}
		}
		for _, topic := range topics {
			if len(embed.Fields) > 19 {
				break
			}
			field := &discordgo.MessageEmbedField{}
			if topic.FirstURL != "" {
				field.Value = topic.FirstURL
			} else {
				field.Value = "[no url]"
			}
			if topic.Text != "" {
				field.Name = topic.Text
				fields = append(fields, field)
			}
		}
		embed.Fields = fields
	}

}

func SearchCommand(ctx *commandclient.Context) {
	channel, err := ctx.Session.State.Channel(ctx.Message.ChannelID)
	isNSFW := false
	if err != nil {
		log.WithError(err).WithField("channel", ctx.Message.ChannelID).Errorln("failed to get channel from state")
	} else if channel != nil {
		isNSFW = channel.NSFW
	}
	answer, err := duckduckgo.Query(client, "\\"+strings.Join(ctx.Args, " "), "ddg-ng/0.1", !isNSFW, true)
	if err != nil {
		log.WithError(err).Errorln("failed to fetch data from ddg")
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("I've encountered an error while trying to access the DDG API: %s", err.Error()))
		return
	}
	response := &discordgo.MessageEmbed{}
	createAnswerEmbed(answer, response)
	if _, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, response); err != nil {
		log.WithError(err).WithField("answer", answer).Errorln("failed to send a ddg answer embed")
	}
}

func createGithubEmbed(repo *github.GithubRepo, embed *discordgo.MessageEmbed) {
	embed.URL = repo.HTMLURL
	embed.Title = repo.Name
	embed.Description = repo.Description
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: repo.Owner.AvatarURL,
	}
	createField := func(name string, value string) *discordgo.MessageEmbedField {
		return &discordgo.MessageEmbedField{
			Name:   name,
			Value:  value,
			Inline: true,
		}
	}
	embed.Fields = []*discordgo.MessageEmbedField{
		createField("Watching:", strconv.Itoa(repo.Watchers)),
		createField("Stars:", fmt.Sprintf("[%d](%s)", repo.StargazersCount, repo.StargazersURL)),
		createField("Issues:", fmt.Sprintf("[%d](%s)", repo.OpenIssuesCount, repo.HTMLURL+"/issues")),
		createField("Commits", repo.HTMLURL+"/commits"),
	}
}

var GITHUB_URL_REGEX = regexp.MustCompile("github\\.com\\/(.+?)\\/([^?\\/\\n]+)")
var GITHUB_REGEX = regexp.MustCompile("(.+?)\\/(.+)")

func GitCommand(ctx *commandclient.Context) {
	args := strings.Join(ctx.Args, " ")
	matches := GITHUB_URL_REGEX.FindStringSubmatch(args)
	if len(matches) < 3 {
		matches = GITHUB_REGEX.FindStringSubmatch(args)
	}
	if len(matches) < 3 {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Please provide a valid github url or use username/reponame")
		return
	}
	username := matches[1]
	reponame := matches[2]
	repo, err := github.Repo(client, username, reponame)
	if err != nil {
		log.WithError(err).Errorln("failed to fetch github data")
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error retrieving repo data: "+err.Error())
		return
	}
	embed := &discordgo.MessageEmbed{}
	createGithubEmbed(repo, embed)
	if _, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed); err != nil {
		log.WithError(err).WithField("repo", repo).Errorln("failed to send a git embed")
	}
}
