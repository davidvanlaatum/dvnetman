package spec

import "dvnetman/api/gen/openapi"

func init() {
	RegisterBuilder(BuildLocation)
}

func BuildLocation(api *OpenAPI) {
	api.AddObject(
		AddObjectOpts{
			Name: "location",
			Tags: []string{"Location"},
			SearchBody: &openapi.Schema{
				Type: "object",
				Properties: map[string]openapi.Schema{
					"ids": {
						Type: "array",
						Items: &openapi.Schema{
							Type:   "string",
							Format: "uuid",
						},
					},
					"name": {
						Type: "string",
					},
					"nameRegex": {
						Type: "string",
					},
					"fields": {
						Type: "array",
						Items: &openapi.Schema{
							Type: "string",
						},
					},
					"parent": {
						Type:   "string",
						Format: "uuid",
					},
					"site": {
						Type:   "string",
						Format: "uuid",
					},
				},
			},
			SearchModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
						"name": {
							Type: "string",
						},
						"description": {
							Type: "string",
						},
						"parent": {
							Ref: objectRef,
						},
						"site": {
							Ref: objectRef,
						},
					},
					Required: []string{"name"},
				},
			),
			GetModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
						"site": {
							Ref: objectRef,
						},
						"parent": {
							Ref: objectRef,
						},
						"name": {
							Type: "string",
						},
						"description": {
							Type: "string",
						},
					},
					Required: []string{"name"},
				},
			),
		},
	)
}
