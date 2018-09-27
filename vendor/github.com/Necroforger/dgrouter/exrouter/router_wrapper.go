package exrouter

import (
	"strings"

	"github.com/Necroforger/dgrouter"
	"github.com/bwmarrin/discordgo"
)

// HandlerFunc ...
type HandlerFunc func(*Context)

// MiddlewareFunc is middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Route wraps dgrouter.Router to use a Context
type Route struct {
	*dgrouter.Route
}

// New returns a new router wrapper
func New() *Route {
	return &Route{
		Route: dgrouter.New(),
	}
}

// On registers a handler function
func (r *Route) On(name string, handler HandlerFunc) *Route {
	return &Route{r.Route.On(name, WrapHandler(handler))}
}

// Group ...
func (r *Route) Group(fn func(rt *Route)) *Route {
	return &Route{r.Route.Group(func(r *dgrouter.Route) {
		fn(&Route{r})
	})}
}

// Use ...
func (r *Route) Use(fn ...MiddlewareFunc) *Route {
	wrapped := []dgrouter.MiddlewareFunc{}
	for _, v := range fn {
		wrapped = append(wrapped, WrapMiddleware(v))
	}
	return &Route{
		r.Route.Use(wrapped...),
	}
}

// WrapMiddleware ...
func WrapMiddleware(mware MiddlewareFunc) dgrouter.MiddlewareFunc {
	return func(next dgrouter.HandlerFunc) dgrouter.HandlerFunc {
		return func(i interface{}) {
			WrapHandler(mware(UnwrapHandler(next)))(i)
		}
	}
}

// OnMatch registers a route with the given matcher
func (r *Route) OnMatch(name string, matcher func(string) bool, handler HandlerFunc) *Route {
	return &Route{r.Route.OnMatch(name, matcher, WrapHandler(handler))}
}

func mention(id string) string {
	return "<@" + id + ">"
}

func nickMention(id string) string {
	return "<@!" + id + ">"
}

// FindAndExecute is a helper method for calling routes
// it creates a context from a message, finds its route, and executes the handler
// it looks for a message prefix which is either the prefix specified or the message is prefixed
// with a bot mention
//    s            : discordgo session to pass to context
//    prefix       : prefix you want the bot to respond to
//    botID        : user ID of the bot to allow you to substitute the bot ID for a prefix
//    m            : discord message to pass to context
func (r *Route) FindAndExecute(s *discordgo.Session, prefix string, botID string, m *discordgo.Message) error {
	var pf string

	// If the message content is only a bot mention and the mention route is not nil, send the mention route
	if r.Default != nil && m.Content == mention(botID) || m.Content == nickMention(botID) {
		r.Default.Handler(NewContext(s, m, []string{""}, r.Default))
		return nil
	}

	// Append a space to the mentions
	bmention := mention(botID) + " "
	nmention := nickMention(botID) + " "

	p := func(t string) bool {
		return strings.HasPrefix(m.Content, t)
	}

	switch {
	case prefix != "" && p(prefix):
		pf = prefix
	case p(bmention):
		pf = bmention
	case p(nmention):
		pf = nmention
	default:
		return dgrouter.ErrCouldNotFindRoute
	}

	command := strings.TrimPrefix(m.Content, pf)
	args := ParseArgs(command)

	if rt, depth := r.FindFull(args...); depth > 0 {
		args = append([]string{strings.Join(args[:depth], string(separator))}, args[depth:]...)
		rt.Handler(NewContext(s, m, args, rt))
	} else {
		return dgrouter.ErrCouldNotFindRoute
	}

	return nil
}

// WrapHandler wraps a dgrouter.HandlerFunc
func WrapHandler(fn HandlerFunc) dgrouter.HandlerFunc {
	if fn == nil {
		return nil
	}
	return func(i interface{}) {
		fn(i.(*Context))
	}
}

// UnwrapHandler unwraps a handler
func UnwrapHandler(fn dgrouter.HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		fn(ctx)
	}
}
