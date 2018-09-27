package dgrouter

import (
	"errors"
)

// Error variables
var (
	ErrCouldNotFindRoute  = errors.New("Could not find route")
	ErrRouteAlreadyExists = errors.New("route already exists")
)

// HandlerFunc is a command handler
type HandlerFunc func(interface{})

// MiddlewareFunc is a middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Group allows you to do things like more easily manage categories
// For example, setting the routes category in the callback will cause
// All future added routes to inherit the category.
// example:
// Group(func (r *Route) {
//    r.Cat("stuff")
//    r.On("thing", nil).Desc("the category of this function will be stuff")
// })
func (r *Route) Group(fn func(r *Route)) *Route {
	rt := New()
	fn(rt)
	for _, v := range rt.Routes {
		r.AddRoute(v)
	}
	return r
}

// Use adds the given middleware func to this route's middleware chain
func (r *Route) Use(fn ...MiddlewareFunc) *Route {
	r.Middleware = append(r.Middleware, fn...)
	return r
}

// On registers a route with the name you supply
//    name    : name of the route to create
//    handler : handler function
func (r *Route) On(name string, handler HandlerFunc) *Route {
	rt := r.OnMatch(name, nil, handler)
	rt.Matcher = NewNameMatcher(rt)
	return rt
}

// OnMatch adds a handler for the given route
//    name    : name of the route to add
//    matcher : matcher function used to match the route
//    handler : handler function for the route
func (r *Route) OnMatch(name string, matcher func(string) bool, handler HandlerFunc) *Route {
	if rt := r.Find(name); rt != nil {
		return rt
	}

	nhandler := handler

	// Add middleware to the handler
	for _, v := range r.Middleware {
		nhandler = v(nhandler)
	}

	rt := &Route{
		Name:     name,
		Category: r.Category,
		Handler:  nhandler,
		Matcher:  matcher,
	}

	r.AddRoute(rt)
	return rt
}

// AddRoute adds a route to the router
// Will return RouteAlreadyExists error on failure
//    route : route to add
func (r *Route) AddRoute(route *Route) error {
	// Check if the route already exists
	if rt := r.Find(route.Name); rt != nil {
		return ErrRouteAlreadyExists
	}

	route.Parent = r
	r.Routes = append(r.Routes, route)
	return nil
}

// RemoveRoute removes a route from the router
//     route : route to remove
func (r *Route) RemoveRoute(route *Route) error {
	for i, v := range r.Routes {
		if v == route {
			r.Routes = append(r.Routes[:i], r.Routes[i+1:]...)
			return nil
		}
	}
	return ErrCouldNotFindRoute
}

// Find finds a route with the given name
// It will return nil if nothing is found
//    name : name of route to find
func (r *Route) Find(name string) *Route {
	for _, v := range r.Routes {
		if v.Matcher(name) {
			return v
		}
	}
	return nil
}

// FindFull a full path of routes by searching through their subroutes
// Until the deepest match is found.
// It will return the route matched and the depth it was found at
//     args : path of route you wish to find
//            ex. FindFull(command, subroute1, subroute2, nonexistent)
//            will return the deepest found match, which will be subroute2
func (r *Route) FindFull(args ...string) (*Route, int) {
	nr := r
	i := 0
	for _, v := range args {
		if rt := nr.Find(v); rt != nil {
			nr = rt
			i++
		} else {
			break
		}
	}
	return nr, i
}

// New returns a new route
func New() *Route {
	return &Route{
		Routes: []*Route{},
	}
}
