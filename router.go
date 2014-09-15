package dispatcher

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Node ...
type Node struct {
	parent   *Node
	children []*Node
	child    map[interface{}]*Node
	name     string
	key      interface{}
	regexp   *regexp.Regexp
	methods  []string
	handlers map[string][]http.Handler
}

// FindOrInsert ...
func (n *Node) FindOrInsert(node *Node) *Node {
	if result := n.Find(node); nil != result {
		return result
	}

	return n.Insert(node)
}

// Find ...
func (n *Node) Find(node *Node) *Node {
	if child := n.child[node.key]; nil != child {
		return child
	}

	var children = n.children

	if 0 < len(node.name) {
		for i := 0; i < len(children); i++ {
			if child := children[i]; child.name == node.name {
				return child
			}
		}
	}

	return nil
}

// Insert ...
func (n *Node) Insert(node *Node) *Node {
	node.parent = n

	if nil == node.key {
		n.children = append(n.children, node)
	} else {
		n.child[node.key] = node
	}

	return node
}

// Match ...
type Match struct {
	node   *Node
	params map[string]interface{}
}

// RouteInformation ...
type RouteInformation struct {
	keys   map[interface{}]bool
	name   string
	regexp *regexp.Regexp
}

// Router ...
type Router struct {
	root *Node
}

// Route ...
func (r *Router) Route(path string) []*Node {
	return r.Define(strings.Split(path, "/"), r.root)
}

// Define ...
func (r *Router) Define(fragments []string, root *Node) []*Node {
	var (
		fragment = fragments[0]
		info     = r.Parse(fragment)
		name     = info.name
		nodes    = make([]*Node, 0, 0)
	)

	for key := range info.keys {
		var node = &Node{
			name:     name,
			key:      key,
			child:    make(map[interface{}]*Node),
			methods:  make([]string, 0, 0),
			handlers: make(map[string][]http.Handler),
		}

		nodes = append(nodes, node)
	}

	if nil != info.regexp {
		var node = &Node{
			name:     name,
			regexp:   info.regexp,
			child:    make(map[interface{}]*Node),
			methods:  make([]string, 0, 0),
			handlers: make(map[string][]http.Handler),
		}

		nodes = append(nodes, node)
	}

	if 0 == len(nodes) {
		var node = &Node{
			name:     name,
			child:    make(map[interface{}]*Node),
			methods:  make([]string, 0, 0),
			handlers: make(map[string][]http.Handler),
		}

		nodes = append(nodes, node)
	}

	for i := 0; i < len(nodes); i++ {
		nodes[i] = root.FindOrInsert(nodes[i])
	}

	if 0 == len(fragments)-1 {
		return nodes
	}

	for i := 0; i < len(nodes); i++ {
		nodes = r.Define(fragments[1:], nodes[i])
	}

	return nodes
}

// Parse ...
func (r *Router) Parse(fragment string) RouteInformation {
	var info RouteInformation
	info.keys = make(map[interface{}]bool)

	if r.IsValidSlug(fragment) {
		info.keys[fragment] = true
		return info
	} else if r.IsPipeSeparatedSlug(fragment) {
		var slugs = strings.Split(fragment, "|")

		for i := 0; i < len(slugs); i++ {
			info.keys[slugs[i]] = true
		}

		return info
	}

	fragment = regexp.MustCompile(`^\:\w+\b`).ReplaceAllStringFunc(fragment, func(s string) string {
		info.name = s[1:]
		return ""
	})

	if 0 == len(fragment) {
		return info
	}

	if regexp.MustCompile(`^\(.+\)$`).MatchString(fragment) {
		fragment = fragment[1 : len(fragment)-1]

		if r.IsPipeSeparatedSlug(fragment) {
			var slugs = strings.Split(fragment, "|")

			for i := 0; i < len(slugs); i++ {
				info.keys[slugs[i]] = true
			}
		} else {
			info.regexp = regexp.MustCompile(fmt.Sprintf("^(%s)$", fragment))
		}
	}

	return info
}

// MatchURL ...
func (r *Router) MatchURL(url string) *Match {
	var (
		root      = r.root
		fragments = strings.Split(url, "/")
		length    = len(fragments)
		match     = new(Match)
	)

	match.params = make(map[string]interface{})

	var (
		fragment string
		name     string
		node     *Node
		nodes    []*Node
		regexp   *regexp.Regexp
	)

top:
	for 0 < length {
		fragment = fragments[0]

		if 1 <= len(fragments) {
			fragments = fragments[1:]
		}

		length = len(fragments)

		if node = root.child[fragment]; nil != node {
			if name = node.name; 0 < len(name) {
				match.params[name] = fragment
			}

			if 0 == length {
				match.node = node
				return match
			}

			root = node
			goto top
		}

		nodes = root.children

		for i := 0; i < len(nodes); i++ {
			node = nodes[i]

			if regexp = node.regexp; nil == regexp || regexp.MatchString(fragment) {
				if name = node.name; 0 < len(name) {
					match.params[name] = fragment
				}

				if 0 == length {
					match.node = node
					return match
				}

				root = node
				goto top
			}
		}
	}

	return nil
}

// IsValidSlug ...
func (r *Router) IsValidSlug(slug string) bool {
	return slug == "" || regexp.MustCompile(`^[\w\.-]+$`).MatchString(slug)
}

// IsPipeSeparatedSlug ...
func (r *Router) IsPipeSeparatedSlug(slug string) bool {
	return regexp.MustCompile(`^[\w\.\-][\w\.\-\|]+[\w\.\-]$`).MatchString(slug)
}
