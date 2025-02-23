package spec

import (
	"dvnetman/api/gen/openapi"
	"dvnetman/pkg/utils"
)

func addCommonProps(s *openapi.Schema) *openapi.Schema {
	s.Properties["id"] = openapi.Schema{
		Type:     "string",
		Format:   "uuid",
		ReadOnly: true,
	}
	s.Properties["created"] = openapi.Schema{
		Type:     "string",
		Format:   dateTime,
		ReadOnly: true,
	}
	s.Properties["updated"] = openapi.Schema{
		Type:     "string",
		Format:   dateTime,
		ReadOnly: true,
	}
	s.Properties["tags"] = openapi.Schema{
		Type: "array",
		Items: &openapi.Schema{
			Ref: "#/components/schemas/Tag",
		},
	}
	s.Properties["version"] = openapi.Schema{
		Type: "integer",
	}
	s.Required = append(s.Required, "id", "version")
	return s
}

type AddObjectOpts struct {
	Name         string
	SearchModal  *openapi.Schema
	SearchParams []openapi.Parameter
	GetModal     *openapi.Schema
	Additional   map[string]*openapi.Schema
	Tags         []string
}

func (o *OpenAPI) AddObject(opts AddObjectOpts) {
	var basePath = "/api/v1/" + opts.Name
	var idPath = basePath + "/{id}"
	resultSchema := o.AddSchema(utils.UCFirst(opts.Name)+"Result", opts.SearchModal)
	resultsSchema := o.AddSchema(
		utils.UCFirst(opts.Name)+"SearchResults", &openapi.Schema{
			Type: "object",
			Properties: map[string]openapi.Schema{
				"items": {
					Type: "array",
					Items: &openapi.Schema{
						Ref: resultSchema,
					},
				},
				"count": {
					Type: "integer",
				},
				"next": {
					Type: "boolean",
				},
			},
			Required: []string{"items", "count", "next"},
		},
	)
	modalSchema := o.AddSchema(utils.UCFirst(opts.Name), opts.GetModal)
	o.AddEndpoint(
		AddEndpointOpts{
			Method:    "get",
			Path:      basePath,
			Operation: "list" + utils.UCFirst(opts.Name) + "s",
			Tags:      opts.Tags,
			Parameters: append(
				[]openapi.Parameter{
					{
						Name:        "page",
						In:          "query",
						Required:    false,
						Description: "Page number",
						Schema: &openapi.Schema{
							Type: "integer",
						},
					},
					{
						Name:        "per_page",
						In:          "query",
						Required:    false,
						Description: "Number of items per page",
						Schema: &openapi.Schema{
							Type: "integer",
						},
					},
					{
						Name:        "sort",
						In:          "query",
						Required:    false,
						Description: "Sort order",
						Schema: &openapi.Schema{
							Type: "string",
						},
					},
				},
				opts.SearchParams...,
			),
			Responses: map[string]openapi.Response{
				"200": {
					Description: "List of " + opts.Name,
					Content: map[string]openapi.Content{
						json: {
							Schema: openapi.Schema{
								Ref: resultsSchema,
							},
						},
					},
				},
			},
		},
	)
	o.AddEndpoint(
		AddEndpointOpts{
			Method:    "post",
			Path:      basePath,
			Operation: "create" + utils.UCFirst(opts.Name),
			Tags:      opts.Tags,
			RequestBody: &openapi.RequestBody{
				Content: map[string]openapi.Content{
					json: {
						Schema: openapi.Schema{
							Ref: modalSchema,
						},
					},
				},
				Required: true,
			},
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Create " + opts.Name,
					Content: map[string]openapi.Content{
						json: {
							Schema: openapi.Schema{
								Ref: modalSchema,
							},
						},
					},
				},
			},
		},
	)
	o.AddEndpoint(
		AddEndpointOpts{
			Path:      idPath,
			Method:    "get",
			Operation: "get" + utils.UCFirst(opts.Name),
			Tags:      opts.Tags,
			Parameters: []openapi.Parameter{
				{
					Name:     "id",
					In:       "path",
					Required: true,
					Schema: &openapi.Schema{
						Type:   "string",
						Format: "uuid",
					},
				},
				{
					Name: "If-None-Match",
					In:   "header",
					Schema: &openapi.Schema{
						Type: "string",
					},
				},
				{
					Name: "If-Modified-Since",
					In:   "header",
					Schema: &openapi.Schema{
						Type:   "string",
						Format: dateTime,
					},
				},
			},
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Get " + opts.Name,
					Content: map[string]openapi.Content{
						json: {
							Schema: openapi.Schema{
								Ref: modalSchema,
							},
						},
					},
				},
				"404": notFoundResponse,
			},
		},
	)
	o.AddEndpoint(
		AddEndpointOpts{
			Method:    "put",
			Path:      idPath,
			Operation: "update" + utils.UCFirst(opts.Name),
			Tags:      opts.Tags,
			Parameters: []openapi.Parameter{
				{
					Name:     "id",
					In:       "path",
					Required: true,
					Schema: &openapi.Schema{
						Type:   "string",
						Format: "uuid",
					},
				},
			},
			RequestBody: &openapi.RequestBody{
				Content: map[string]openapi.Content{
					json: {
						Schema: openapi.Schema{
							Ref: modalSchema,
						},
					},
				},
				Required: true,
			},
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Update " + opts.Name,
					Content: map[string]openapi.Content{
						json: {
							Schema: openapi.Schema{
								Ref: modalSchema,
							},
						},
					},
				},
			},
		},
	)
	o.AddEndpoint(
		AddEndpointOpts{
			Method:    "delete",
			Path:      idPath,
			Operation: "delete" + utils.UCFirst(opts.Name),
			Tags:      opts.Tags,
			Parameters: []openapi.Parameter{
				{
					Name:     "id",
					In:       "path",
					Required: true,
					Schema: &openapi.Schema{
						Type:   "string",
						Format: "uuid",
					},
				},
			},
			Responses: map[string]openapi.Response{
				"204": {
					Description: "Delete " + opts.Name,
				},
			},
		},
	)
	o.Components.Schemas["Stats"].Properties[opts.Name+"Count"] = openapi.Schema{
		Type: "integer",
	}
	for k, v := range opts.Additional {
		o.AddSchema(k, v)
	}
}
