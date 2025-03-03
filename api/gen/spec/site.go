package spec

import "dvnetman/api/gen/openapi"

func init() {
	RegisterBuilder(BuildSite)
}

func BuildSite(api *OpenAPI) {
	api.AddObject(
		AddObjectOpts{
			Name: "site",
			Tags: []string{"Site"},
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
					},
					Required: []string{"name"},
				},
			),
			GetModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
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
