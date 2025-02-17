package ymlutil

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

func Walk(y *yaml.Node, f func(*yaml.Node, []string), path ...string) {
	switch y.Kind {
	case yaml.DocumentNode:
		Walk(y.Content[0], f)
	case yaml.MappingNode:
		for i := 0; i < len(y.Content); i += 2 {
			Walk(y.Content[i+1], f, append(path, y.Content[i].Value)...)
		}
	case yaml.ScalarNode:
		f(y, path)
	case yaml.SequenceNode:
		for i, n := range y.Content {
			Walk(n, f, append(path, fmt.Sprintf("[%d]", i))...)
		}
	default:
		panic(fmt.Sprintf("unexpected node kind: %v", y.Kind))
	}
}
