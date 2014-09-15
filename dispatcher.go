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

// Namespace ...
type Namespace struct {
	prefix     []string
	dispatcher *Dispatcher
}

// Get ...
func (n *Namespace) Get(url string, handlers ...http.Handler) *Namespace {
	url = fmt.Sprintf("%s/%s", strings.Join(n.prefix, "/"), strings.TrimLeft(url, "/"))
	n.dispatcher.Route(methods.Get, url, handlers...)
	return n
}

// Put ...
func (n *Namespace) Put(url string, handlers ...http.Handler) *Namespace {
	url = fmt.Sprintf("%s/%s", strings.Join(n.prefix, "/"), strings.TrimLeft(url, "/"))
	n.dispatcher.Route(methods.Put, url, handlers...)
	return n
}

// Post ...
func (n *Namespace) Post(url string, handlers ...http.Handler) *Namespace {
	url = fmt.Sprintf("%s/%s", strings.Join(n.prefix, "/"), strings.TrimLeft(url, "/"))
	n.dispatcher.Route(methods.Post, url, handlers...)
	return n
}

// Patch ...
func (n *Namespace) Patch(url string, handlers ...http.Handler) *Namespace {
	url = fmt.Sprintf("%s/%s", strings.Join(n.prefix, "/"), strings.TrimLeft(url, "/"))
	n.dispatcher.Route(methods.Patch, url, handlers...)
	return n
}

// Delete ...
func (n *Namespace) Delete(url string, handlers ...http.Handler) *Namespace {
	url = fmt.Sprintf("%s/%s", strings.Join(n.prefix, "/"), strings.TrimLeft(url, "/"))
	n.dispatcher.Route(methods.Delete, url, handlers...)
	return n
}

// Namespace ...
func (n *Namespace) Namespace(prefix string) *Namespace {
	prefix = strings.TrimSuffix(strings.TrimPrefix(prefix, "/"), "/")
	return &Namespace{prefix: append(n.prefix, prefix), dispatcher: n.dispatcher}
}

// Dispatcher ...
type Dispatcher struct {
	router     *Router
	middleware []http.Handler
}

// Namespace ...
func (d *Dispatcher) Namespace(prefix string) *Namespace {
	return &Namespace{prefix: []string{strings.TrimRight(prefix, "/")}, dispatcher: d}
}

// Route ...
func (d *Dispatcher) Route(method, url string, handlers ...http.Handler) *Dispatcher {
	method = strings.ToUpper(method)

	var nodes = d.router.Route(url)

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
	return d.Route(methods.Get, url, handlers...)
}

// Put ...
func (d *Dispatcher) Put(url string, handlers ...http.Handler) *Dispatcher {
	return d.Route(methods.Put, url, handlers...)
}

// Post ...
func (d *Dispatcher) Post(url string, handlers ...http.Handler) *Dispatcher {
	return d.Route(methods.Post, url, handlers...)
}

// Patch ...
func (d *Dispatcher) Patch(url string, handlers ...http.Handler) *Dispatcher {
	return d.Route(methods.Patch, url, handlers...)
}

// Delete ...
func (d *Dispatcher) Delete(url string, handlers ...http.Handler) *Dispatcher {
	return d.Route(methods.Delete, url, handlers...)
}

// ServeHTTP ...
func (d *Dispatcher) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	var (
		match  = d.router.MatchURL(request.URL.Path)
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
		router     = new(Router)
		root       = new(Node)
	)

	root.child = make(map[interface{}]*Node)
	router.root = root
	dispatcher.router = router
	dispatcher.middleware = make([]http.Handler, 0, 0)
	return dispatcher
}
