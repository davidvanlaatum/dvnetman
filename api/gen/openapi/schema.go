package openapi

import (
	"dvnetman/pkg/utils"
	"gopkg.in/yaml.v3"
	"strings"
)

type Info struct {
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
	Version     string `yaml:"version,omitempty"`
}

type Server struct {
	URL string `yaml:"url,omitempty"`
}

type Schema struct {
	Type       string            `yaml:"type,omitempty"`
	Items      *Schema           `yaml:"items,omitempty"`
	Ref        string            `yaml:"$ref,omitempty"`
	Properties map[string]Schema `yaml:"properties,omitempty"`
	Required   []string          `yaml:"required,omitempty"`
	Format     string            `yaml:"format,omitempty"`
	ReadOnly   bool              `yaml:"readOnly,omitempty"`
	Enum       []string          `yaml:"enum,omitempty"`
	Example    string            `yaml:"example,omitempty"`
	Pattern    string            `yaml:"pattern,omitempty"`
}

type Parameter struct {
	Name            string  `yaml:"name,omitempty"`
	In              string  `yaml:"in,omitempty"`
	Description     string  `yaml:"description,omitempty"`
	Required        bool    `yaml:"required,omitempty"`
	Schema          *Schema `yaml:"schema,omitempty"`
	AllowEmptyValue bool    `yaml:"allowEmptyValue,omitempty"`
}

type Content struct {
	Schema  Schema      `yaml:"schema,omitempty"`
	Example interface{} `yaml:"example,omitempty"`
}

type Response struct {
	Description string             `yaml:"description,omitempty"`
	Content     map[string]Content `yaml:"content,omitempty"`
}

type Endpoint struct {
	OperationID string                `yaml:"operationId,omitempty"`
	Parameters  []Parameter           `yaml:"parameters,omitempty"`
	Responses   map[string]Response   `yaml:"responses,omitempty"`
	RequestBody *RequestBody          `yaml:"requestBody,omitempty"`
	Tags        []string              `yaml:"tags,omitempty"`
	Security    []map[string][]string `yaml:"security,omitempty"`
}

type Components struct {
	Schemas map[string]*Schema `yaml:"schemas,omitempty"`
}

type Paths map[string]map[string]Endpoint

func pathLess(a, b string) bool {
	ap := strings.Split(a, "/")
	bp := strings.Split(b, "/")
	for x := 0; x < len(ap) && x < len(bp); x++ {
		if ap[x] != bp[x] {
			if ap[x][0] == '{' && bp[x][0] != '{' {
				return false
			}
			if ap[x][0] != '{' && bp[x][0] == '{' {
				return true
			}
			return ap[x] < bp[x]
		}
	}
	return a < b
}

func (o Paths) Sorted() (values []utils.SortedMapEntries[string, []utils.SortedMapEntries[string, Endpoint]]) {
	return utils.MapTo(
		utils.MapSortedByKey(
			o, func(a, b utils.SortedMapEntries[string, map[string]Endpoint]) bool {
				return pathLess(a.Key, b.Key)
			},
		), func(
			v utils.SortedMapEntries[string, map[string]Endpoint],
		) utils.SortedMapEntries[string, []utils.SortedMapEntries[string, Endpoint]] {
			return utils.SortedMapEntries[string, []utils.SortedMapEntries[string, Endpoint]]{
				Key:   v.Key,
				Value: utils.MapSortedByKey(v.Value, utils.MapSortedByKeyString),
			}
		},
	)
}

func subDocumentToNode(v interface{}) (*yaml.Node, error) {
	n := &yaml.Node{}
	if b, err := yaml.Marshal(v); err != nil {
		return nil, err
	} else if err = yaml.Unmarshal(b, n); err != nil {
		return nil, err
	}
	return n.Content[0], nil
}

func (o Paths) MarshalYAML() (interface{}, error) {
	pathsNode := &yaml.Node{Kind: yaml.MappingNode}
	for _, path := range o.Sorted() {
		pathsNode.Content = append(pathsNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: path.Key})
		pathNode := &yaml.Node{Kind: yaml.MappingNode}
		pathsNode.Content = append(pathsNode.Content, pathNode)
		for _, method := range path.Value {
			pathNode.Content = append(pathNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: method.Key})
			if methodNode, err := subDocumentToNode(method.Value); err != nil {
				return nil, err
			} else {
				pathNode.Content = append(pathNode.Content, methodNode)
			}
		}
	}
	return pathsNode, nil
}

type RequestBody struct {
	Content  map[string]Content `yaml:"content,omitempty"`
	Required bool               `yaml:"required,omitempty"`
}

type OpenAPI struct {
	OpenAPI    string     `yaml:"openapi,omitempty"`
	Info       Info       `yaml:"info,omitempty"`
	Servers    []Server   `yaml:"servers,omitempty"`
	Paths      Paths      `yaml:"paths,omitempty"`
	Components Components `yaml:"components,omitempty"`
}
