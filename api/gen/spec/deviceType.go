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
					"model": {
						Type: "string",
					},
					"modelRegex": {
						Type: "string",
					},
					"manufacturer": {
						Type: "array",
						Items: &openapi.Schema{
							Type:   "string",
							Format: "uuid",
						},
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
						"model": {
							Type: "string",
						},
						"manufacturer": {
							Ref: objectRef,
						},
					},
				},
			),
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
