package openapi

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestSchemaMarshal(t *testing.T) {
	r := require.New(t)
	api := &OpenAPI{
		Info: Info{
			Title:   "DVNetMan",
			Version: "1.0.0",
		},
		Paths: map[string]map[string]Endpoint{
			"/abc/xyz": {
				"get": {
					OperationID: "getAbcXyz",
				},
				"post": {
					OperationID: "postAbcXyz",
				},
			},
			"/users/{id}": {
				"get": {
					OperationID: "listUsers",
				},
			},
			"/users/current": {
				"get": {
					OperationID: "getCurrentUser",
				},
			},
		},
	}
	buf := &bytes.Buffer{}
	e := yaml.NewEncoder(buf)
	e.SetIndent(2)
	err := e.Encode(api)
	r.NoError(err)
	r.Equal(
		`info:
  title: DVNetMan
  version: 1.0.0
paths:
  /abc/xyz:
    get:
      operationId: getAbcXyz
    post:
      operationId: postAbcXyz
  /users/current:
    get:
      operationId: getCurrentUser
  /users/{id}:
    get:
      operationId: listUsers
`, buf.String(),
	)
}
