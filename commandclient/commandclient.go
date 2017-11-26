package commandclient

import (
	// "database/sql"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// Context stores command context
type Context struct {
	Client  *CommandClient
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Args    []string
}

// New creates new command client
func New(prefix string /* , db *sql.DB */) *CommandClient {
	return &CommandClient{
		Prefix:   prefix,
		Register: make(map[string]func(ctx *Context)),
		// DB:       db,
	}
}

// CommandClient session/command  register
type CommandClient struct {
	Prefix   string
	Register map[string]func(ctx *Context)
	// DB       *sql.DB
}

// Parse parses a message event
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
