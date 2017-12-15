package gitplugin

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aurieh/ddg-ng/commandclient"
	"github.com/aurieh/ddg-ng/github"
	"github.com/aurieh/ddg-ng/utils"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strconv"
	"strings"
)

func createGithubEmbed(repo *github.Repo, embed *discordgo.MessageEmbed) {
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

// GithubURLRegex matches github repo urls
var GithubURLRegex = regexp.MustCompile(`github\.com\/(.+?)\/([^?\/\n]+)`)

// GithubRegex matches repo shorthands
var GithubRegex = regexp.MustCompile(`(.+?)\/(.+)`)

// GitCommand displays git repo info
func GitCommand(ctx *commandclient.Context) {
	args := strings.Join(ctx.Args, " ")
	matches := GithubURLRegex.FindStringSubmatch(args)
	if len(matches) < 3 {
		matches = GithubRegex.FindStringSubmatch(args)
	}
	if len(matches) < 3 {
		_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Please provide a valid github url or use username/reponame")
		if err != nil {
			log.WithError(err).Errorln("couldn't send message")
		}
		return
	}
	username := matches[1]
	reponame := matches[2]
	repo, err := github.GetRepo(utils.Client, username, reponame)
	if err != nil {
		log.WithError(err).Errorln("failed to fetch github data")
		_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error retrieving repo data: "+err.Error())
		if err != nil {
			log.WithError(err).Errorln("couldn't send message")
		}
		return
	}
	embed := &discordgo.MessageEmbed{}
	createGithubEmbed(repo, embed)
	if _, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed); err != nil {
		log.WithError(err).WithField("repo", repo).Errorln("failed to send a git embed")
	}
}
