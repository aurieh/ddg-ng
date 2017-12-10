package commandclient

import (
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
func New(prefix string) *CommandClient {
	return &CommandClient{
		Prefix:   prefix,
		Register: make(map[string]func(ctx *Context)),
	}
}

// CommandClient session/command  register
type CommandClient struct {
	Prefix   string
	Register map[string]func(ctx *Context)
	OnMissingPrefix func(s *discordgo.Session, m *discordgo.MessageCreate)
	OnUnknownCommand func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
	OnSuccessfulInvoke func(ctx *Context, command func(ctx *Context))
}

// Parse parses a message event
func (p *CommandClient) Parse(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, p.Prefix) {
		if p.OnMissingPrefix != nil {
			p.OnMissingPrefix(s, m)
		}
		return
	}

	args := strings.Split(strings.TrimSpace(m.Content[len(p.Prefix):]), " ")
	if len(args) < 1 {
		return
	}

	if command, ok := p.Register[args[0]]; ok {
		ctx := &Context{
			Client:  p,
			Session: s,
			Message: m,
			Args:    args[1:],
		}
		command(ctx)
		if p.OnSuccessfulInvoke != nil {
			p.OnSuccessfulInvoke(ctx, command)
		}
	} else if p.OnUnknownCommand != nil {
		p.OnUnknownCommand(s, m, args)
	}
}
