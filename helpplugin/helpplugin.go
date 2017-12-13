package helpplugin

import (
	"fmt"
	"github.com/aurieh/ddg-ng/commandclient"
	"strings"
)

const info = "```tex" + `
$$ ddg-ng $$

# OWNER: taciturasa
# DEVELOPERS: taciturasa, aurieh
# LIB: discordgo by bwmarrin

% ddg-ng is a frontend for the DuckDuckGo Instant Answers API.
% It is made for quick searches and info grabbing on Discord.
% It supports most of ddg answer types and search syntax.
% Ready to !bang?

Powered by DuckDuckGo {{https://duckduckgo.com/about}}
` + "```"

func InfoCommand(ctx *commandclient.Context) {
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, info)
}

const help = `
**All commands are called with ddg!command or @ddg command.**

**ddg! search terms and @ddg search terms will also search DDG.**

Available commands: %s
`
func HelpCommand(ctx *commandclient.Context) {
	var commands []string
	for name := range ctx.Client.Register {
		commands = append(commands, "`"+name+"`")
	}
	message := fmt.Sprintf(help, strings.Join(commands, ", "))
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, message)
}

const server = "_For more help on this bot, please visit: <https://discord.gg/011iDaqaFcbzbEsMz>_"
func ServerCommand(ctx *commandclient.Context) {
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, server)
}

const addbot = "_To add this bot to your server, use this link https://thats-a.link/e9bccc_"
func AddbotCommand(ctx *commandclient.Context) {
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, addbot)
}
