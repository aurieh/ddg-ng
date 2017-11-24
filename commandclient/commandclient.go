package commandclient

import (
	// "database/sql"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Context struct {
	Client  *CommandClient
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Args    []string
}

func New(prefix string/* , db *sql.DB */) *CommandClient {
	return &CommandClient{
		Prefix:   prefix,
		Register: make(map[string]func(ctx *Context)),
		// DB:       db,
	}
}

type CommandClient struct {
	Prefix   string
	Register map[string]func(ctx *Context)
	// DB       *sql.DB
}

func (p *CommandClient) Parse(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, p.Prefix) {
		return
	}

	args := strings.Split(strings.TrimSpace(m.Content[len(p.Prefix):]), " ")
	if len(args) < 1 {
		return
	}

	if command, ok := p.Register[args[0]]; ok {
		command(&Context{
			Client:  p,
			Session: s,
			Message: m,
			Args:    args[1:],
		})
	}
}
