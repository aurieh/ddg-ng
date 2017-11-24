package main

import (
	// "database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/aurieh/ddg-ng/commandclient"
	"github.com/aurieh/ddg-ng/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	// _ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	pflag.String("token", "", "discord token")
	pflag.Bool("debug", false, "debug level")
	// pflag.String("db", "", "pgsql connection string")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvPrefix("DDG")

	viper.BindEnv("token")
	viper.BindEnv("debug")

	// viper.SetDefault("db", "")
	// viper.BindEnv("db")

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
	// if !viper.IsSet("db") {
	// 	log.Fatalln("no db conn string specified")
	// }
	// db, err := sql.Open("postgres", viper.GetString("db"))
	// if err != nil {
	// 	log.WithError(err).Fatalln("failed while connecting to the database")
	// }

	if !viper.IsSet("token") {
		log.Fatalln("no token specified")
	}
	dg, err := discordgo.New("Bot " + viper.GetString("token"))
	client := commandclient.New("ddg!")
	client.Register["stats"] = commands.StatsCommand
	client.Register["search"] = commands.SearchCommand
	client.Register["git"] = commands.GitCommand

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