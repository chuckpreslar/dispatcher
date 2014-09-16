package dispatcher

import "testing"

import "github.com/stretchr/testify/assert"

func TestStringMatch(t *testing.T) {
	var (
		trie = new(Trie)
		root = new(Node)
	)

	root.child = make(map[interface{}]*Node)
	trie.root = root

	var (
		nodes = trie.Establish("/api/v1/users")
		match = trie.MatchURL("/api/v1/users")
	)

	assert.NotNil(t, nodes, "failed to return nodes from `Establish`")
	assert.NotNil(t, match, "failed to return match from `MatchURL`")

	var node = match.node

	assert.Equal(t, node.key, "users", "failed to return correct node")
}

func TestPipedStringMatch(t *testing.T) {
	var (
		trie = new(Trie)
		root = new(Node)
	)

	root.child = make(map[interface{}]*Node)
	trie.root = root

	var (
		nodes  = trie.Establish("/api/v1/users|posts")
		umatch = trie.MatchURL("/api/v1/users")
		pmatch = trie.MatchURL("/api/v1/posts")
	)

	assert.NotNil(t, nodes, "failed to return nodes from `Establish`")
	assert.NotNil(t, umatch, "failed to return match from `MatchURL`")
	assert.NotNil(t, pmatch, "failed to return match from `MatchURL`")

	var (
		unode   = umatch.node
		pnode   = pmatch.node
		uparent = unode.parent
		pparent = pnode.parent
	)

	assert.Equal(t, unode.key, "users", "failed to return appropriate named node")
	assert.Equal(t, pnode.key, "posts", "failed to return appropriate named node")

	assert.Equal(t, uparent, pparent, "failed to return appropriate parent")
	assert.Equal(t, uparent.key, "v1", "failed to return appropriate named node")
	assert.Equal(t, pparent.key, "v1", "failed to return appropriate named node")
}

func TestWildcardMatch(t *testing.T) {
	var (
		trie = new(Trie)
		root = new(Node)
	)

	root.child = make(map[interface{}]*Node)
	trie.root = root

	var (
		nodes = trie.Establish("/api/v1/users/:id")
		match = trie.MatchURL("/api/v1/users/1")
	)

	assert.NotNil(t, nodes, "failed to return nodes from `Establish`")
	assert.NotNil(t, match, "failed to return match from `MatchURL`")
	assert.Equal(t, match.params["id"], "1", "failed to return correct node")

	var node = match.node

	assert.Equal(t, node.name, "id", "failed to return correct node")
}

func TestRegexpMatch(t *testing.T) {
	var (
		trie = new(Trie)
		root = new(Node)
	)

	root.child = make(map[interface{}]*Node)
	trie.root = root

	var (
		nodes = trie.Establish(`/api/v1/users/([0-9]{1})`)
		match = trie.MatchURL("/api/v1/users/1")
	)

	assert.NotNil(t, nodes, "failed to return nodes from `Establish`")
	assert.NotNil(t, match, "failed to return match from `MatchURL`")
}

func TestNamedRegexpMatch(t *testing.T) {
	var (
		trie = new(Trie)
		root = new(Node)
	)

	root.child = make(map[interface{}]*Node)
	trie.root = root

	var (
		nodes = trie.Establish(`/api/v1/users/:id([0-9]{1})`)
		match = trie.MatchURL("/api/v1/users/1")
	)

	assert.NotNil(t, nodes, "failed to return nodes from `Establish`")
	assert.NotNil(t, match, "failed to return match from `MatchURL`")
	assert.Equal(t, match.params["id"], "1", "failed to return correct node")
	var node = match.node

	assert.Equal(t, node.name, "id", "failed to return correct node")

}
