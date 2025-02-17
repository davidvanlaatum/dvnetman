package ymlutil

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestWalk(t *testing.T) {
	r := require.New(t)
	x := &yaml.Node{}
	err := yaml.Unmarshal(
		[]byte(`servers:
  - url: xyz
`), x,
	)
	r.NoError(err)

	var paths [][]string
	Walk(
		x, func(n *yaml.Node, path []string) {
			paths = append(paths, path)
			t.Logf("%s: %s", path, n.Value)
		},
	)
	r.Equal([][]string{{"servers", "[0]", "url"}}, paths)
}
