package dgrouter

// Route is a command route
type Route struct {
	// Routes is a slice of subroutes
	Routes []*Route

	Name        string
	Aliases     []string
	Description string
	Category    string

	// Matcher is a function that determines
	// If this route will be matched
	Matcher func(string) bool

	// Handler is the Handler for this route
	Handler HandlerFunc

	// Default route for responding to bot mentions
	Default *Route

	// The parent for this route
	Parent *Route

	// Middleware to be applied when adding subroutes
	Middleware []MiddlewareFunc
}

// Desc sets this routes description
func (r *Route) Desc(description string) *Route {
	r.Description = description
	return r
}

// Cat sets this route's category
func (r *Route) Cat(category string) *Route {
	r.Category = category
	return r
}

// Alias appends aliases to this route's alias list
func (r *Route) Alias(aliases ...string) *Route {
	r.Aliases = append(r.Aliases, aliases...)
	return r
}
