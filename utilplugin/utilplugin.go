package utilplugin

import (
	"github.com/aurieh/ddg-ng/commandclient"
	"github.com/aurieh/ddg-ng/stats"
	log "github.com/Sirupsen/logrus"
)

// StatsCommand get bot stats
func StatsCommand(ctx *commandclient.Context) {
	statstring := stats.GetStatsString()
	if _, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, statstring); err != nil {
		log.WithError(err).WithField("stats", statstring).Errorln("failed to send stats")
	}
}
