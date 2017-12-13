package stats

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"runtime"
	"strconv"
	"time"
)

var startup = time.Now()

// GetStatsString get a formatted stats string
func GetStatsString(s *discordgo.Session) string {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)

	table.SetHeader([]string{"Name", "Value"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	table.Append([]string{"Go Runtime", runtime.Version()})
	table.Append([]string{"Goroutines", strconv.Itoa(runtime.NumGoroutine())})
	table.Append([]string{"CPU Threads", strconv.Itoa(runtime.NumCPU())})
	table.Append([]string{"Memory Allocated", humanize.Bytes(stats.Alloc)})
	table.Append([]string{"Memory Total", humanize.Bytes(stats.Sys)})
	table.Append([]string{"Memory GCed", humanize.Bytes(stats.TotalAlloc)})
	table.Append([]string{"Guilds", strconv.Itoa(len(s.State.Guilds))})
	table.Append([]string{"Up Since", humanize.Time(startup)})

	fmt.Fprintln(buf, "```")
	table.Render()
	fmt.Fprint(buf, "\n```")
	return buf.String()
}
