package utilplugin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aurieh/ddg-ng/commandclient"
	"github.com/aurieh/ddg-ng/stats"
)

// StatsCommand get bot stats
func StatsCommand(ctx *commandclient.Context) {
	statstring := stats.GetStatsString(ctx.Session)
	if _, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, statstring); err != nil {
		log.WithError(err).WithField("stats", statstring).Errorln("failed to send stats")
	}
}
