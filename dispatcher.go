package dispatcher

import (
	"fmt"
	"net/http"
	"strings"
)

import (
	"github.com/chuckpreslar/dispatcher/headers"
	"github.com/chuckpreslar/dispatcher/methods"
	"github.com/chuckpreslar/dispatcher/statuses"
)

// Route ...
type Route struct {
	url        string
	dispatcher *Dispatcher
}

// Get ...
func (r *Route) Get(handlers ...http.Handler) *Route {
	r.dispatcher.Get(r.url, handlers...)
	return r
}

// Put ...
func (r *Route) Put(handlers ...http.Handler) *Route {
	r.dispatcher.Put(r.url, handlers...)
	return r
}

// Post ...
func (r *Route) Post(handlers ...http.Handler) *Route {
	r.dispatcher.Post(r.url, handlers...)
	return r
}

// Patch ...
func (r *Route) Patch(handlers ...http.Handler) *Route {
	r.dispatcher.Patch(r.url, handlers...)
	return r
}

// Delete ...
func (r *Route) Delete(handlers ...http.Handler) *Route {
	r.dispatcher.Delete(r.url, handlers...)
	return r
}

// Dispatcher ...
type Dispatcher struct {
	trie       *Trie
	middleware []http.Handler
}

// RegisterRouteHandlers ...
func (d *Dispatcher) RegisterRouteHandlers(method, url string, handlers ...http.Handler) *Dispatcher {
	method = strings.ToUpper(method)

	var nodes = d.trie.Establish(url)

	for i := 0; i < len(nodes); i++ {
		var node = nodes[i]

		node.handlers[method] = append(node.handlers[method], handlers...)

		if -1 == strings.Index(strings.Join(node.methods, ""), method) {
			node.methods = append(node.methods, method)
		}
	}

	return d
}

// Get ...
func (d *Dispatcher) Get(url string, handlers ...http.Handler) *Dispatcher {
	return d.RegisterRouteHandlers(methods.Get, url, handlers...)
}

// Put ...
func (d *Dispatcher) Put(url string, handlers ...http.Handler) *Dispatcher {
	return d.RegisterRouteHandlers(methods.Put, url, handlers...)
}

// Post ...
func (d *Dispatcher) Post(url string, handlers ...http.Handler) *Dispatcher {
	return d.RegisterRouteHandlers(methods.Post, url, handlers...)
}

// Patch ...
func (d *Dispatcher) Patch(url string, handlers ...http.Handler) *Dispatcher {
	return d.RegisterRouteHandlers(methods.Patch, url, handlers...)
}

// Delete ...
func (d *Dispatcher) Delete(url string, handlers ...http.Handler) *Dispatcher {
	return d.RegisterRouteHandlers(methods.Delete, url, handlers...)
}

// Route ...
func (d *Dispatcher) Route(url string) *Route {
	return &Route{url: url, dispatcher: d}
}

// ServeHTTP ...
func (d *Dispatcher) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	var (
		match  = d.trie.MatchURL(request.URL.Path)
		header = response.Header()
	)

	if nil == match {
		return
	}

	var (
		node     = match.node
		method   = strings.ToUpper(request.Method)
		handlers = append(node.handlers[method], d.middleware...)
	)

	if methods.Options == method {
		header.Set(headers.Allow, strings.Join(node.methods, ","))
		response.WriteHeader(statuses.OK)
		return
	}

	for i := 0; i < len(handlers); i++ {
		handlers[i].ServeHTTP(response, request)
	}
}

// Use ...
func (d *Dispatcher) Use(handlers ...http.Handler) *Dispatcher {
	d.middleware = append(d.middleware, handlers...)
	return d
}

// Listen ...
func (d *Dispatcher) Listen(port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), d)
}

// New ...
func New() *Dispatcher {
	var (
		dispatcher = new(Dispatcher)
		trie       = new(Trie)
		root       = new(Node)
	)

	root.child = make(map[interface{}]*Node)
	trie.root = root
	dispatcher.trie = trie
	dispatcher.middleware = make([]http.Handler, 0, 0)
	return dispatcher
}
