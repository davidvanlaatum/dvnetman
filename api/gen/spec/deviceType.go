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
