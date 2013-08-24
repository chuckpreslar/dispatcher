package dispatcher

import (
  "testing"
)

type reader struct{}
func (r reader) Read(p []byte) (n int, err error) { return }

func TestNewRoute(t *testing.T) {
  route := NewRoute("/test/:required/:optional?", false)
  path := "/test/one/two"

  if !route.matcher.MatchString(path) {
    t.Error("Expected route to match path with required and optional parameters supplied.")
  } else if path = "/test/one"; !route.matcher.MatchString(path) {
    t.Error("Expected route to match path with optional parameter missing.")
  }
}
