package spec

import "dvnetman/api/gen/openapi"

func init() {
	RegisterBuilder(BuildManufacturer)
}

func BuildManufacturer(api *OpenAPI) {
	api.AddObject(
		AddObjectOpts{
			Name: "manufacturer",
			Tags: []string{"Manufacturer"},
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
					},
					Required: []string{"id", "name"},
				},
			),
			GetModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
						"name": {
							Type: "string",
						},
					},
					Required: []string{"id", "name"},
				},
			),
		},
	)
}
