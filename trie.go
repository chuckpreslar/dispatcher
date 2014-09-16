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

// Trie ...
type Trie struct {
	root *Node
}

// Establish ...
func (t *Trie) Establish(path string) []*Node {
	return t.Define(strings.Split(path, "/"), t.root)
}

// Define ...
func (t *Trie) Define(fragments []string, root *Node) []*Node {
	var (
		fragment = fragments[0]
		info     = t.Parse(fragment)
		name     = info.name
		nodes    = make([]*Node, 0, 0)
		temp     = make([]*Node, 0, 0)
	)

	if 0 <= len(fragments)-1 {
		fragments = fragments[1:]
	}

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

	if 0 == len(fragments) {
		return nodes
	}

	for i := 0; i < len(nodes); i++ {
		if 0 < len(fragments) {
			temp = t.Define(fragments, nodes[i])
		}
	}

	return temp
}

// Parse ...
func (t *Trie) Parse(fragment string) RouteInformation {
	var info RouteInformation
	info.keys = make(map[interface{}]bool)

	if t.IsValidSlug(fragment) {
		info.keys[fragment] = true
		return info
	} else if t.IsPipeSeparatedSlug(fragment) {
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

		if t.IsPipeSeparatedSlug(fragment) {
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
func (t *Trie) MatchURL(url string) *Match {
	var (
		root      = t.root
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
func (t *Trie) IsValidSlug(slug string) bool {
	return slug == "" || regexp.MustCompile(`^[\w\.-]+$`).MatchString(slug)
}

// IsPipeSeparatedSlug ...
func (t *Trie) IsPipeSeparatedSlug(slug string) bool {
	return regexp.MustCompile(`^[\w\.\-][\w\.\-\|]+[\w\.\-]$`).MatchString(slug)
}
