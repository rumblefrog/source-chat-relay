package exrouter

import (
	"fmt"
	"sync"

	"github.com/Necroforger/dgrouter"

	"github.com/bwmarrin/discordgo"
)

// Context represents a command context
type Context struct {
	// Route is the route that this command came from
	Route *dgrouter.Route
	Msg   *discordgo.Message
	Ses   *discordgo.Session

	// List of arguments supplied with the command
	Args Args

	// Vars that can be optionally set using the Set and Get functions
	vmu  sync.RWMutex
	Vars map[string]interface{}
}

// Set sets a variable on the context
func (c *Context) Set(key string, d interface{}) {
	c.vmu.Lock()
	c.Vars[key] = d
	c.vmu.Unlock()
}

// Get retrieves a variable from the context
func (c *Context) Get(key string) interface{} {
	if c, ok := c.Vars[key]; ok {
		return c
	}
	return nil
}

// Reply replies to the sender with the given message
func (c *Context) Reply(args ...interface{}) (*discordgo.Message, error) {
	return c.Ses.ChannelMessageSend(c.Msg.ChannelID, fmt.Sprint(args...))
}

// ReplyEmbed replies to the sender with an embed
func (c *Context) ReplyEmbed(args ...interface{}) (*discordgo.Message, error) {
	return c.Ses.ChannelMessageSendEmbed(c.Msg.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprint(args...),
	})
}

// Guild retrieves a guild from the state or restapi
func (c *Context) Guild(guildID string) (*discordgo.Guild, error) {
	g, err := c.Ses.State.Guild(guildID)
	if err != nil {
		g, err = c.Ses.Guild(guildID)
	}
	return g, err
}

// Channel retrieves a channel from the state or restapi
func (c *Context) Channel(channelID string) (*discordgo.Channel, error) {
	ch, err := c.Ses.State.Channel(channelID)
	if err != nil {
		ch, err = c.Ses.Channel(channelID)
	}
	return ch, err
}

// Member retrieves a member from the state or restapi
func (c *Context) Member(guildID, userID string) (*discordgo.Member, error) {
	m, err := c.Ses.State.Member(guildID, userID)
	if err != nil {
		m, err = c.Ses.GuildMember(guildID, userID)
	}
	return m, err
}

// NewContext returns a new context from a message
func NewContext(s *discordgo.Session, m *discordgo.Message, args Args, route *dgrouter.Route) *Context {
	return &Context{
		Route: route,
		Msg:   m,
		Ses:   s,
		Args:  args,
		Vars:  map[string]interface{}{},
	}
}
