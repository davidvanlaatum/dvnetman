package spec

import (
	"dvnetman/api/gen/openapi"
	"dvnetman/pkg/file"
	"dvnetman/pkg/logger"
	"dvnetman/version"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type OpenAPI struct {
	openapi.OpenAPI
	log logger.Logger
}

const objectRef = "#/components/schemas/ObjectReference"
const errorRef = "#/components/schemas/APIErrorModal"
const json = "application/json"
const dateTime = "date-time"

var badRequestResponse = openapi.Response{
	Description: "Bad request",
	Content: map[string]openapi.Content{
		json: {
			Schema: openapi.Schema{
				Type: "object",
				Items: &openapi.Schema{
					Ref: errorRef,
				},
			},
			Example: map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"code":    "bad_request",
						"message": "Bad request",
					},
				},
			},
		},
	},
}
var notFoundResponse = openapi.Response{
	Description: "Not found",
	Content: map[string]openapi.Content{
		json: {
			Schema: openapi.Schema{
				Type: "object",
				Items: &openapi.Schema{
					Ref: errorRef,
				},
			},
			Example: map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"code":    "not_found",
						"message": "Not found",
					},
				},
			},
		},
	},
}
var internalServerErrorResponse = openapi.Response{
	Description: "Internal server error",
	Content: map[string]openapi.Content{
		json: {
			Schema: openapi.Schema{
				Type: "object",
				Items: &openapi.Schema{
					Ref: errorRef,
				},
			},
			Example: map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"code":    "internal_server_error",
						"message": "Internal server error",
					},
				},
			},
		},
	},
}

func (o *OpenAPI) AddSchema(name string, s *openapi.Schema) string {
	if _, found := o.Components.Schemas[name]; found {
		panic("schema " + name + " already exists")
	}
	o.Components.Schemas[name] = s
	return "#/components/schemas/" + name
}

func (o *OpenAPI) WriteOpenAPISpec(path string) (err error) {
	o.log.Info().Key("path", path).Msg("Writing OpenAPI spec")
	var f file.FileUpdate
	if f, err = file.NewFileUpdate(path, nil); err != nil {
		return
	}
	defer f.Close()
	e := yaml.NewEncoder(f)
	defer e.Close()
	e.SetIndent(2)
	if err = e.Encode(o.OpenAPI); err != nil {
		return errors.Wrap(err, "failed to encode yaml")
	}
	return
}

var builders []func(*OpenAPI)

func RegisterBuilder(f func(*OpenAPI)) {
	builders = append(builders, f)
}

func NewSpec(log logger.Logger) *OpenAPI {
	api := &OpenAPI{
		OpenAPI: openapi.OpenAPI{
			OpenAPI: "3.1.0",
			Info: openapi.Info{
				Title:       "DVNetMan",
				Description: "DVNetMan",
				Version:     version.Version,
			},
			Servers: []openapi.Server{
				{
					URL: "/",
				},
			},
			Paths: map[string]map[string]openapi.Endpoint{},
			Components: openapi.Components{
				Schemas: map[string]*openapi.Schema{},
			},
		},
		log: log,
	}
	errorMsg := api.AddSchema(
		"ErrorMessage", &openapi.Schema{
			Type: "object",
			Properties: map[string]openapi.Schema{
				"code": {
					Type: "string",
				},
				"message": {
					Type: "string",
				},
			},
			Required: []string{"code", "message"},
		},
	)
	api.AddSchema(
		"APIErrorModal", &openapi.Schema{
			Type: "object",
			Properties: map[string]openapi.Schema{
				"errors": {
					Type: "array",
					Items: &openapi.Schema{
						Ref: errorMsg,
					},
				},
			},
		},
	)
	api.AddSchema(
		"ObjectReference", &openapi.Schema{
			Type: "object",
			Properties: map[string]openapi.Schema{
				"id": {
					Type:   "string",
					Format: "uuid",
				},
				"displayName": {
					Type:     "string",
					ReadOnly: true,
				},
			},
			Required: []string{"id"},
		},
	)
	api.AddSchema(
		"Tag", &openapi.Schema{
			Type: "object",
			Properties: map[string]openapi.Schema{
				"id": {
					Type:   "string",
					Format: "uuid",
				},
				"name": {
					Type: "string",
				},
			},
		},
	)
	api.AddSchema(
		"Stats", &openapi.Schema{
			Type:       "object",
			Properties: map[string]openapi.Schema{},
		},
	)

	api.AddEndpoint(
		AddEndpointOpts{
			Method:    "get",
			Path:      "/api/v1/stats",
			Operation: "GetStats",
			Tags:      []string{"Stats"},
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Stats",
					Content: map[string]openapi.Content{
						json: {
							Schema: openapi.Schema{
								Ref: "#/components/schemas/Stats",
							},
						},
					},
				},
			},
		},
	)
	for _, builder := range builders {
		builder(api)
	}
	return api
}
