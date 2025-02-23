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
