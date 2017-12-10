package main

import (
	"github.com/aurieh/ddg-ng/commandclient"
	"github.com/aurieh/ddg-ng/gitplugin"
	"github.com/aurieh/ddg-ng/searchplugin"
	"github.com/aurieh/ddg-ng/utilplugin"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
	"syscall"
	log "github.com/Sirupsen/logrus"
)

func init() {
	pflag.String("token", "", "discord token")
	pflag.Bool("debug", false, "debug level")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvPrefix("DDG")

	viper.BindEnv("token")
	viper.BindEnv("debug")

	viper.SetConfigType("yaml")
	viper.SetConfigName("ddg")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/ddg")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.WithError(err).Fatalln("failed while reading config")
		}
	}

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	if !viper.IsSet("token") {
		log.Fatalln("no token specified")
	}
	dg, err := discordgo.New("Bot " + viper.GetString("token"))
	client := commandclient.New("ddg!")
	client.OnMissingPrefix = onMissingPrefix
	client.OnUnknownCommand = onUnknownCommand
	client.Register["stats"] = utilplugin.StatsCommand
	client.Register["search"] = searchplugin.SearchCommand
	client.Register["git"] = gitplugin.GitCommand

	if err != nil {
		log.WithError(err).Fatalln("failed while creating a discord instance")
	}
	dg.AddHandler(client.Parse)

	err = dg.Open()
	if err != nil {
		log.Fatalln(err)
	}
	log.Infoln("Discord connection open. Press C^c or send SIGTERM to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Infoln("Bye")
	err = dg.Close()
	if err != nil {
		log.Errorln(err)
	}
}

func onMissingPrefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.WithError(err).WithField("channelID", m.ChannelID).Errorln("failed to get channel from state")
		return
	}
	if channel.Type != discordgo.ChannelTypeDM {
		return
	}
	searchplugin.SendQueryResult(s, m, m.Content)
}

func onUnknownCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	searchplugin.SendQueryResult(s, m, strings.Join(args, " "))
}
