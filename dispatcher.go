package dispatcher

import (
	"fmt"
	"net/http"
	"strings"
)

import (
	"github.com/chuckpreslar/dispatcher/methods"
	"github.com/chuckpreslar/dispatcher/statuses"
)

// Dispatcher ...
type Dispatcher struct {
	router *Router
}

// Route ...
func (d *Dispatcher) Route(method, url string, handler http.Handler) *Dispatcher {
	method = strings.ToUpper(method)
	var nodes = d.router.Route(url)

	for i := 0; i < len(nodes); i++ {
		var node = nodes[i]

		node.handlers[method] = append(node.handlers[method], handler)

		if -1 == strings.Index(strings.Join(node.methods, ""), method) {
			node.methods = append(node.methods, method)
		}
	}

	return d
}

// Get ...
func (d *Dispatcher) Get(url string, handler http.Handler) *Dispatcher {
	return d.Route(methods.Get, url, handler)
}

// Put ...
func (d *Dispatcher) Put(url string, handler http.Handler) *Dispatcher {
	return d.Route(methods.Put, url, handler)
}

// Post ...
func (d *Dispatcher) Post(url string, handler http.Handler) *Dispatcher {
	return d.Route(methods.Post, url, handler)
}

// Patch ...
func (d *Dispatcher) Patch(url string, handler http.Handler) *Dispatcher {
	return d.Route(methods.Patch, url, handler)
}

// Delete ...
func (d *Dispatcher) Delete(url string, handler http.Handler) *Dispatcher {
	return d.Route(methods.Delete, url, handler)
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
		handlers = node.handlers[method]
	)

	if methods.Options == method {
		header.Set("Allow", strings.Join(node.methods, ","))
		response.WriteHeader(statuses.OK)
		return
	}

	for i := 0; i < len(handlers); i++ {
		handlers[i].ServeHTTP(response, request)
	}
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

	return dispatcher
}
