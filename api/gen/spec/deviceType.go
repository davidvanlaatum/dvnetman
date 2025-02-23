package spec

import "dvnetman/api/gen/openapi"

func init() {
	RegisterBuilder(BuildDeviceType)
}

func BuildDeviceType(api *OpenAPI) {
	api.AddObject(
		AddObjectOpts{
			Name: "deviceType",
			Tags: []string{"Device"},
			SearchModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
						"model": {
							Type: "string",
						},
						"manufacturer": {
							Ref: objectRef,
						},
					},
				},
			),
			SearchParams: []openapi.Parameter{
				{
					Name:        "ids",
					In:          "query",
					Description: "IDs of the device types",
					Schema: &openapi.Schema{
						Type: "array",
						Items: &openapi.Schema{
							Type:   "string",
							Format: "uuid",
						},
					},
				},
				{
					Name:        "model",
					In:          "query",
					Description: "Model of the device type",
					Schema: &openapi.Schema{
						Type: "string",
					},
				},
				{
					Name:        "modelRegex",
					In:          "query",
					Description: "Model of the device type",
					Schema: &openapi.Schema{
						Type: "string",
					},
				},
				{
					Name:        "fields",
					In:          "query",
					Description: "Fields to return",
					Schema: &openapi.Schema{
						Type: "array",
						Items: &openapi.Schema{
							Type: "string",
						},
					},
				},
			},
			GetModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
						"model": {
							Type: "string",
						},
						"manufacturer": {
							Ref: objectRef,
						},
					},
				},
			),
		},
	)
}
